package snapshots

import (
	"context"
	"fmt"
	"time"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/helpers"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/kubestorageclasses"

	"k8s.io/client-go/dynamic"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/uuid"
	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func RestoreSnapshot(ctx context.Context, snapshotName string, sourceNamespace string, targetName string, targetNamespace string, restoreStorageClass string) error {
	kubeconfig, err := k8s.GetClientConfig()
	if err != nil {
		return err
	}
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	k8sclient, err := k8s.GetClientWithConfig(config)
	if err != nil {
		return err
	}
	if targetNamespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		targetNamespace = contextNamespace
	}
	if sourceNamespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		sourceNamespace = contextNamespace
	}

	if restoreStorageClass == "" {
		restoreStorageClass, err = kubestorageclasses.GetDefaultStorageClassName(ctx, k8sclient)
		if err != nil {
			return err
		}
	}

	// make sure the storage class can be used for restore purposes
	foundStorageClass, err := k8sclient.StorageV1().StorageClasses().Get(ctx, restoreStorageClass, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to retrieve storage class: %w", err)
	}
	if foundStorageClass.VolumeBindingMode == nil {
		return fmt.Errorf("failed to detect volume binding mode of storage class: requires Immediate")
	}
	if *foundStorageClass.VolumeBindingMode != storagev1.VolumeBindingImmediate {
		return fmt.Errorf("storage class volume binding mode is not immedidate")
	}

	snapshotUnstructued, err := client.Resource(volumesnapshotResource).Namespace(sourceNamespace).Get(ctx, snapshotName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// convert to the snapshot object
	snapshot := volumesnapshotv1.VolumeSnapshot{}
	runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotUnstructued.Object, &snapshot)

	// Check that the snapshot is ready to use
	if snapshot.Status == nil || snapshot.Status.ReadyToUse == nil || !*snapshot.Status.ReadyToUse {
		return fmt.Errorf("snapshot is not ready for use")
	}

	fmt.Printf("using snapshot %s for restoring\n", snapshot.GetName())
	restorePVCName := fmt.Sprintf("%s-%s", targetName, uuid.NewString())
	// get size from the volumesnapshot restoresize
	storageSize := *snapshot.Status.RestoreSize

	sourcePVC := ""
	if snapshot.Spec.Source.PersistentVolumeClaimName != nil {
		sourcePVC = *snapshot.Spec.Source.PersistentVolumeClaimName
	}
	_, err = k8sclient.CoreV1().PersistentVolumeClaims(sourceNamespace).Create(ctx, &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      restorePVCName,
			Namespace: sourceNamespace,
			Labels: map[string]string{
				"acloud-toolkit.k8s.avisi.cloud/snapshot-reference": string(snapshot.GetUID()),
				"acloud-toolkit.k8s.avisi.cloud/target-pvc":         targetName,
				"acloud-toolkit.k8s.avisi.cloud/source-snapshot":    snapshotName,
				"acloud-toolkit.k8s.avisi.cloud/source-pvc":         sourcePVC,
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			StorageClassName: helpers.String(restoreStorageClass),
			DataSource: &apiv1.TypedLocalObjectReference{
				APIGroup: helpers.String("snapshot.storage.k8s.io"),
				Kind:     "VolumeSnapshot",
				Name:     snapshotName,
			},
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"storage": storageSize,
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("created PVC %s...\n", restorePVCName)

	// wait until PVC has a persistent volume
	pvc, err := k8sclient.CoreV1().PersistentVolumeClaims(sourceNamespace).Get(ctx, restorePVCName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for {
		if pvc.Spec.VolumeName != "" {
			break
		}
		time.Sleep(1 * time.Second)

		pvc, err = k8sclient.CoreV1().PersistentVolumeClaims(sourceNamespace).Get(ctx, restorePVCName, metav1.GetOptions{})
		if err != nil {
			return err
		}
	}
	fmt.Printf("PVC has volume %s...\n", pvc.Spec.VolumeName)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()
	pvc, err = k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctxWithTimeout, k8sclient, restorePVCName, sourceNamespace)
	if err != nil {
		return err
	}

	err = helpers.RetryWithCancel(ctxWithTimeout, 3, 2*time.Second, func() error {
		return k8s.SetPVReclaimPolicyToRetain(ctx, k8sclient, pvc)
	})
	if err != nil {
		return err
	}

	// Delete the PVC
	err = k8sclient.CoreV1().PersistentVolumeClaims(sourceNamespace).Delete(ctx, restorePVCName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("deleted the PVC %s...\n", restorePVCName)

	err = helpers.RetryWithCancel(ctxWithTimeout, 3, 2*time.Second, func() error {
		err = k8s.RemoveClaimRefOfPV(ctxWithTimeout, k8sclient, pvc)
		if err != nil {
			return err
		}
		claimRef := v1.ObjectReference{Name: targetName, Namespace: targetNamespace}
		return k8s.SetClaimRefOfPV(ctxWithTimeout, k8sclient, pvc.Spec.VolumeName, claimRef)
	})
	if err != nil {
		return err
	}

	// Create a new PVC in the target namespace.
	_, err = k8sclient.CoreV1().PersistentVolumeClaims(targetNamespace).Create(ctx, &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      targetName,
			Namespace: targetNamespace,
			Labels: map[string]string{
				"acloud-toolkit.k8s.avisi.cloud/restored":           "true",
				"acloud-toolkit.k8s.avisi.cloud/snapshot-reference": string(snapshot.GetUID()),
				"acloud-toolkit.k8s.avisi.cloud/target-pvc":         fmt.Sprintf("%s", targetName),
				"acloud-toolkit.k8s.avisi.cloud/source-snapshot":    fmt.Sprintf("%s", snapshotName),
				"acloud-toolkit.k8s.avisi.cloud/source-pvc":         fmt.Sprintf("%s", sourcePVC),
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			StorageClassName: helpers.String(restoreStorageClass),
			AccessModes:      []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"storage": storageSize,
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("created a new PVC %s in namespace %s...\n", targetName, targetNamespace)

	// validate that the new PVC has the correct persistent volume claimed.
	targetPVC, err := k8sclient.CoreV1().PersistentVolumeClaims(targetNamespace).Get(ctx, targetName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for {
		if targetPVC.Spec.VolumeName != "" {
			break
		}
		time.Sleep(1 * time.Second)

		targetPVC, err = k8sclient.CoreV1().PersistentVolumeClaims(targetNamespace).Get(ctx, targetName, metav1.GetOptions{})
		if err != nil {
			return err
		}
	}
	if targetPVC.Spec.VolumeName != pvc.Spec.VolumeName {
		fmt.Printf("warning: the restored pvc does not have the expected volume name. Expected: %q, found %q\n", pvc.Spec.VolumeName, targetPVC.Spec.VolumeName)
	}
	fmt.Printf("restore completed\n")
	// TODO: patch PV if it was originaly marked as PersistentVolumeReclaimPolicy = Delete, back to Delete

	return nil
}
