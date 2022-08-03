package snapshots

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/helpers"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"

	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func SnapshotCreate(snapshotName string, targetNamespace string, targetName string, snapshotClassName string) error {
	kubeconfig, err := k8s.GetClientCmd()
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

	// wait until PVC has a persistent volume
	pvc, err := k8sclient.CoreV1().PersistentVolumeClaims(targetNamespace).Get(context.TODO(), targetName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("PVC with volume %s found...\n", pvc.Spec.VolumeName)

	volumesnapshotRes := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"}

	var volumeSnapshotClassName *string
	if strings.TrimSpace(snapshotClassName) != "" {
		volumeSnapshotClassName = helpers.String(strings.TrimSpace(snapshotClassName))
	}
	// convert to the snapshot object
	createSnapshot := volumesnapshotv1.VolumeSnapshot{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "snapshot.storage.k8s.io/v1",
			Kind:       "VolumeSnapshot",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: targetNamespace,
		},
		Spec: volumesnapshotv1.VolumeSnapshotSpec{
			Source: volumesnapshotv1.VolumeSnapshotSource{
				PersistentVolumeClaimName: helpers.String(targetName),
			},
			VolumeSnapshotClassName: volumeSnapshotClassName,
		},
	}
	snapshotUnstructued, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&createSnapshot)
	if err != nil {
		return err
	}

	fmt.Printf("Creating snapshot %q for PVC %q...\n", snapshotName, targetName)
	_, err = client.Resource(volumesnapshotRes).Namespace(targetNamespace).Create(context.TODO(), &unstructured.Unstructured{
		Object: snapshotUnstructued,
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot %q created for PVC %q...\n", snapshotName, targetName)
	fmt.Printf("Waiting until %q is ready for use...\n", snapshotName)
	for {
		snapshotUnstructued, err := client.Resource(volumesnapshotRes).Namespace(targetNamespace).Get(context.TODO(), snapshotName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// convert to the snapshot object
		snapshot := volumesnapshotv1.VolumeSnapshot{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotUnstructued.Object, &snapshot)

		// Check that the snapshot is ready to use
		if snapshot.Status != nil && snapshot.Status.ReadyToUse != nil && *snapshot.Status.ReadyToUse {
			fmt.Printf("Snapshot %q is ready for use...\n", snapshotName)
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Snapshot %q completed...\n", snapshotName)

	return nil
}
