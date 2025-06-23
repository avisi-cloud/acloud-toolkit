package snapshots

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/utils/ptr"
)

func ImportSnapshotFromRawID(ctx context.Context, snapshotName, targetNamespace, snapshotClassName string, rawID string) error {
	kubeconfig, err := k8s.GetClientConfig()
	if err != nil {
		return fmt.Errorf("error getting kubeconfig: %w", err)
	}
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("error loading kubeconfig client: %w", err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error loading dynamic client: %w", err)
	}
	if targetNamespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return fmt.Errorf("error checking namespace: %w", err)
		}
		targetNamespace = contextNamespace
	}

	snapshotContentName := fmt.Sprintf("%s-%s", snapshotName, uuid.NewString())

	volumeSnapshotClassName := ""
	if strings.TrimSpace(snapshotClassName) != "" {
		// TODO: add validation for snapshot classname
		volumeSnapshotClassName = strings.TrimSpace(snapshotClassName)
	}

	volumeSnapshotClass, err := geVolumeSnapshotClassOrDefault(ctx, client, volumeSnapshotClassName)
	if err != nil {
		return fmt.Errorf("no volume snapshot class found: %w", err)
	}

	// Restore the volumesnapshot content
	snapshotContent := volumesnapshotv1.VolumeSnapshotContent{
		TypeMeta: metav1.TypeMeta{
			APIVersion: volumesnapshotv1.SchemeGroupVersion.String(),
			Kind:       "VolumeSnapshotContent",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: snapshotContentName,
			Labels: map[string]string{
				"k8s.avisi.cloud/snapshot-import": "true",
			},
		},
		Spec: volumesnapshotv1.VolumeSnapshotContentSpec{
			VolumeSnapshotClassName: &volumeSnapshotClass.Name,
			// make sure the deletion policy is retain - snapshot is imported and may not be managed by this cluster
			DeletionPolicy: volumesnapshotv1.VolumeSnapshotContentRetain,
			Driver:         volumeSnapshotClass.Driver,
			Source: volumesnapshotv1.VolumeSnapshotContentSource{
				SnapshotHandle: ptr.To(rawID),
			},
			VolumeSnapshotRef: v1.ObjectReference{
				APIVersion: volumesnapshotv1.SchemeGroupVersion.String(),
				Kind:       "VolumeSnapshot",
				Name:       snapshotName,
				Namespace:  targetNamespace,
			},
		},
	}

	snapshotContentUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&snapshotContent)
	if err != nil {
		return fmt.Errorf("error setting up snapshotContent unstructured resource: %w", err)
	}
	fmt.Printf("Creating snapshotcontent %q..\n", snapshotContentName)
	_, err = client.Resource(volumesnapshotContentResource).Create(ctx, &unstructured.Unstructured{
		Object: snapshotContentUnstructured,
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating volumesnapshotcontent: %w", err)
	}

	// convert to the snapshot object
	createSnapshot := volumesnapshotv1.VolumeSnapshot{
		TypeMeta: metav1.TypeMeta{
			APIVersion: volumesnapshotv1.SchemeGroupVersion.String(),
			Kind:       "VolumeSnapshot",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: targetNamespace,
			Labels: map[string]string{
				"k8s.avisi.cloud/snapshot-import": "true",
			},
		},
		Spec: volumesnapshotv1.VolumeSnapshotSpec{
			Source: volumesnapshotv1.VolumeSnapshotSource{
				VolumeSnapshotContentName: ptr.To(snapshotContentName),
			},
			VolumeSnapshotClassName: &volumeSnapshotClass.Name,
		},
	}
	snapshotUnstructued, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&createSnapshot)
	if err != nil {
		return fmt.Errorf("error setting up snapshot unstructured resource: %w", err)
	}

	fmt.Printf("Creating snapshot %q..\n", snapshotName)
	_, err = client.Resource(volumesnapshotResource).Namespace(targetNamespace).Create(ctx, &unstructured.Unstructured{
		Object: snapshotUnstructued,
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating volumesnapshot: %w", err)
	}

	fmt.Printf("Completed importing snapshot %q into CSI snapshot %s..\n", rawID, snapshotName)
	return nil
}
