package snapshots

import (
	"context"
	"fmt"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ErrNoVolumeSnapshotClassFound   = fmt.Errorf("no default volume snapshot class found")
	ErrNoVolumeSnapshotContentFound = fmt.Errorf("no default volume snapshot content found")
)

var (
	volumesnapshotResource        = schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"}
	volumesnapshotContentResource = schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshotcontents"}
	volumesnapshotClassResource   = schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshotclasses"}
)

func geVolumeSnapshotClassOrDefault(ctx context.Context, client *dynamic.DynamicClient, volumeSnapshotClassName string) (*volumesnapshotv1.VolumeSnapshotClass, error) {
	if volumeSnapshotClassName != "" {
		volumeSnapshotClass, err := geVolumeSnapshotClass(ctx, client, volumeSnapshotClassName)
		if err != nil {
			return nil, fmt.Errorf("volume snapshot class not found: %w", err)
		}
		return volumeSnapshotClass, nil
	}

	volumeSnapshotClass, err := getDefaultVolumeSnapshotClass(ctx, client)
	if err != nil && err != ErrNoVolumeSnapshotClassFound {
		return nil, fmt.Errorf("failed to find default volumesnapshotclass: %w", err)
	}
	return volumeSnapshotClass, nil
}

func geVolumeSnapshotClass(ctx context.Context, client *dynamic.DynamicClient, className string) (*volumesnapshotv1.VolumeSnapshotClass, error) {
	snapshotClass, err := client.Resource(volumesnapshotClassResource).Get(ctx, className, metav1.GetOptions{})
	if err != nil && !kubeerrors.IsNotFound(err) {
		return nil, err
	}

	if snapshotClass != nil {
		volumeSnapshotClass := volumesnapshotv1.VolumeSnapshotClass{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotClass.UnstructuredContent(), &volumeSnapshotClass); err != nil {
			return nil, fmt.Errorf("failed to convert: %w", err)
		}
		return &volumeSnapshotClass, nil
	}
	return nil, ErrNoVolumeSnapshotClassFound
}

func getDefaultVolumeSnapshotClass(ctx context.Context, client *dynamic.DynamicClient) (*volumesnapshotv1.VolumeSnapshotClass, error) {
	snapshotClassNames, err := client.Resource(volumesnapshotClassResource).List(ctx, metav1.ListOptions{})
	if err != nil && !kubeerrors.IsNotFound(err) {
		return nil, err
	}

	defaultSnapshotClass := volumesnapshotv1.VolumeSnapshotClass{}
	for _, vsc := range snapshotClassNames.Items {
		annotations := vsc.GetAnnotations()

		value, ok := annotations["snapshot.storage.kubernetes.io/is-default-class"]
		if ok && value == "true" {

			runtime.DefaultUnstructuredConverter.FromUnstructured(vsc.UnstructuredContent(), &defaultSnapshotClass)
			return &defaultSnapshotClass, nil
		}
	}
	return nil, ErrNoVolumeSnapshotClassFound
}

func geVolumeSnapshotContent(ctx context.Context, client *dynamic.DynamicClient, name string) (*volumesnapshotv1.VolumeSnapshotContent, error) {
	volumeSnapshotContent, err := client.Resource(volumesnapshotContentResource).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !kubeerrors.IsNotFound(err) {
		return nil, err
	}

	if volumeSnapshotContent != nil {
		snapshotContentData := volumesnapshotv1.VolumeSnapshotContent{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(volumeSnapshotContent.UnstructuredContent(), &snapshotContentData); err != nil {
			return nil, fmt.Errorf("failed to convert: %w", err)
		}
		return &snapshotContentData, nil
	}
	return nil, ErrNoVolumeSnapshotContentFound
}

// resolveNamespaces returns the namespaces to operate on: all namespaces in the cluster when allNamespaces is set,
// otherwise the given namespace, falling back to the namespace of the current Kubernetes context.
func resolveNamespaces(ctx context.Context, kubeconfig clientcmd.ClientConfig, k8sclient kubernetes.Interface, namespace string, allNamespaces bool) ([]string, error) {
	if allNamespaces {
		namespaceList, err := k8sclient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		namespaces := make([]string, 0, len(namespaceList.Items))
		for _, ns := range namespaceList.Items {
			namespaces = append(namespaces, ns.GetName())
		}
		return namespaces, nil
	}
	if namespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return nil, err
		}
		namespace = contextNamespace
	}
	return []string{namespace}, nil
}

// listVolumeSnapshots returns the VolumeSnapshots in the given namespaces matching labelSelector (empty matches all).
func listVolumeSnapshots(ctx context.Context, client dynamic.Interface, namespaces []string, labelSelector string) ([]volumesnapshotv1.VolumeSnapshot, error) {
	snapshots := []volumesnapshotv1.VolumeSnapshot{}
	for _, namespace := range namespaces {
		snapshotList, err := client.Resource(volumesnapshotResource).Namespace(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return nil, err
		}
		if snapshotList == nil {
			continue
		}
		for _, item := range snapshotList.Items {
			snapshot := volumesnapshotv1.VolumeSnapshot{}
			if err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &snapshot); err != nil {
				return nil, fmt.Errorf("failed to convert volumesnapshot %s/%s: %w", item.GetNamespace(), item.GetName(), err)
			}
			snapshots = append(snapshots, snapshot)
		}
	}
	return snapshots, nil
}
