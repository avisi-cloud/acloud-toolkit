package k8s

import (
	"context"
	"fmt"
	"time"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/helpers"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewVolumeMount(name, path string, readOnly bool) *v1.VolumeMount {
	return &v1.VolumeMount{
		Name:      name,
		MountPath: path,
		ReadOnly:  readOnly,
	}
}

func NewPersistentVolumeClaimVolume(name, claimName string, readOnly bool) *v1.Volume {
	return &v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: claimName,
				ReadOnly:  readOnly,
			},
		},
	}
}

func SetPVReclaimPolicyToRetain(ctx context.Context, k8sClient kubernetes.Interface, pvc *v1.PersistentVolumeClaim) error {
	// Get the persistent volume, ensure it's set to Retain.
	pv, err := k8sClient.CoreV1().PersistentVolumes().Get(ctx, pvc.Spec.VolumeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volume %q: %w", pvc.Name, err)
	}

	if pv.Spec.PersistentVolumeReclaimPolicy != v1.PersistentVolumeReclaimRetain {
		fmt.Printf("PV %s does not have retain as the reclaim policy, updating ...\n", pvc.Spec.VolumeName)
		pv.Spec.PersistentVolumeReclaimPolicy = v1.PersistentVolumeReclaimRetain
		_, err = k8sClient.CoreV1().PersistentVolumes().Update(ctx, pv, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update reclaim policy persistent volume %q: %w", pvc.Name, err)
		}

		// give kube some time to catch up
		time.Sleep(1 * time.Second)
	} else {
		fmt.Printf("PV %s already has retain as the reclaim policy...\n", pvc.Spec.VolumeName)
	}

	// give kube some time to catch up
	time.Sleep(1 * time.Second)
	return nil
}

func GetPersistentVolumeClaimAndCheckForVolumes(ctx context.Context, k8sClient kubernetes.Interface, pvcName string, namespace string) (*v1.PersistentVolumeClaim, error) {
	pvc, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get persistent volume claim %q: %w", pvcName, err)
	}
	for {
		if pvc.Spec.VolumeName != "" {
			break
		}

		pvc, err = k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get persistent volume claim %q: %w", pvcName, err)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return pvc, err
}

func RemoveClaimRefOfPV(ctx context.Context, k8sClient kubernetes.Interface, pvc *v1.PersistentVolumeClaim) error {
	// Update the PVC to remove the claimRef
	pv, err := k8sClient.CoreV1().PersistentVolumes().Get(ctx, pvc.Spec.VolumeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	pv.Spec.ClaimRef = nil
	_, err = k8sClient.CoreV1().PersistentVolumes().Update(ctx, pv, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove claimref of persistent volume claim %q: %w", pv.Name, err)
	}
	fmt.Printf("removed the PV %s claim ref to %s...\n", pvc.Spec.VolumeName, pvc.Name)
	return nil
}

func SetClaimRefOfPV(ctx context.Context, k8sClient kubernetes.Interface, volumeName string, claimRef v1.ObjectReference) error {
	// Update the PVC to remove the claimRef
	pv, err := k8sClient.CoreV1().PersistentVolumes().Get(ctx, volumeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volume %q: %w", volumeName, err)
	}
	pv.Spec.ClaimRef = &claimRef
	_, err = k8sClient.CoreV1().PersistentVolumes().Update(ctx, pv, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update claimref persistent volume claim %q: %w", volumeName, err)
	}
	fmt.Printf("set the PV %s claim ref to %s in namespace %s...\n", volumeName, claimRef.Name, claimRef.Namespace)
	return nil
}

func CreatePersistentVolumeClaim(ctx context.Context, k8sClient kubernetes.Interface, pvcName, namespace, storageClass string, storageSize resource.Quantity) error {
	_, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			StorageClassName: helpers.String(storageClass),
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": storageSize,
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create persistent volume claim %q: %w", pvcName, err)
	}

	fmt.Printf("created a new PVC %s in namespace %s...\n", pvcName, namespace)
	return nil
}
