package snapshots

import (
	"context"
	"fmt"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	"k8s.io/client-go/dynamic"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
	"github.com/avisi-cloud/acloud-toolkit/pkg/table"
)

// List prints the VolumeSnapshots in the namespace (or all namespaces). When batchID is set, only snapshots of that
// batch (labeled SnapshotBatchLabelKey=batchID) are listed.
func List(ctx context.Context, namespace string, allNamespaces, fetchSnapshotHandle bool, batchID string) error {
	labelSelector := ""
	if batchID != "" {
		if err := ValidateSnapshotBatchID(batchID); err != nil {
			return err
		}
		labelSelector = SnapshotBatchLabelSelector(batchID)
	}
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
	namespace, err = resolveNamespace(kubeconfig, namespace, allNamespaces)
	if err != nil {
		return err
	}

	listSnapshots, err := listVolumeSnapshots(ctx, client, namespace, labelSelector)
	if err != nil {
		return err
	}

	// Format output
	body := make([][]string, 0, len(listSnapshots))
	for _, snapshot := range listSnapshots {
		sourceName := ""
		if snapshot.Spec.Source.PersistentVolumeClaimName != nil {
			sourceName = fmt.Sprintf("pvc/%s", *snapshot.Spec.Source.PersistentVolumeClaimName)
		}
		if snapshot.Spec.Source.VolumeSnapshotContentName != nil {
			sourceName = fmt.Sprintf("volumesnapshot/%s", *snapshot.Spec.Source.VolumeSnapshotContentName)
		}
		contentName := getVolumeSnapshotContentNameForSnapshot(snapshot)

		snapshotClassName := ""
		if snapshot.Spec.VolumeSnapshotClassName != nil {
			snapshotClassName = *snapshot.Spec.VolumeSnapshotClassName
		}
		size := ""
		if snapshot.Status != nil && snapshot.Status.RestoreSize != nil {
			size = snapshot.Status.RestoreSize.String()
		}

		snapshotHandle := ""
		if fetchSnapshotHandle {
			snapshotHandle, err = getSnapshotHandle(ctx, contentName, client)
			if err != nil {
				return err
			}
		}

		body = append(body, []string{
			snapshot.GetNamespace(),
			snapshot.Name,
			sourceName,
			contentName,
			size,
			fmt.Sprintf("%v", snapshot.Status != nil && snapshot.Status.ReadyToUse != nil && *snapshot.Status.ReadyToUse),
			snapshotClassName,
			snapshot.GetLabels()[SnapshotBatchLabelKey],
			snapshotHandle,
		})
	}

	table.Print([]string{
		"Namespace",
		"Name",
		"Source",
		"Content",
		"Size",
		"Ready",
		"Classname",
		"Batch",
		"Snapshot Handle",
	}, body)

	return nil
}

func getSnapshotHandle(ctx context.Context, contentName string, client *dynamic.DynamicClient) (string, error) {
	if contentName != "" {
		content, err := geVolumeSnapshotContent(ctx, client, contentName)
		if err != nil && err != ErrNoVolumeSnapshotContentFound {
			return "", err
		}
		if content.Status.SnapshotHandle != nil {
			return *content.Status.SnapshotHandle, nil
		} else if content.Spec.Source.SnapshotHandle != nil {
			return *content.Spec.Source.SnapshotHandle, nil
		}
	}
	return "", nil
}

func getVolumeSnapshotContentNameForSnapshot(snapshot volumesnapshotv1.VolumeSnapshot) string {
	contentName := ""
	if snapshot.Spec.Source.VolumeSnapshotContentName != nil {
		contentName = *snapshot.Spec.Source.VolumeSnapshotContentName
	}
	if snapshot.Status != nil && snapshot.Status.BoundVolumeSnapshotContentName != nil {
		contentName = *snapshot.Status.BoundVolumeSnapshotContentName
	}
	return contentName
}
