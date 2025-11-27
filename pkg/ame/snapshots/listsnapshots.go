package snapshots

import (
	"context"
	"fmt"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
	"github.com/avisi-cloud/acloud-toolkit/pkg/table"
)

func List(ctx context.Context, namespace string, allNamespaces, fetchSnapshotHandle bool) error {
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
	if namespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		namespace = contextNamespace
	}

	// Collect namespaces, either based a single (default) namespace, or collect all namespaces in the cluster
	// and use those to query the snapshots in cluster.
	namespaces := []string{}
	if allNamespaces {
		k8sclient := k8s.GetClientOrDie()
		namespaceList, err := k8sclient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, ns := range namespaceList.Items {
			namespaces = append(namespaces, ns.GetName())
		}
	} else {
		namespaces = append(namespaces, namespace)
	}

	// collect all listSnapshots for each namespace
	listSnapshots := []volumesnapshotv1.VolumeSnapshot{}
	for _, namespace := range namespaces {
		snapshotUnstructured, err := client.Resource(volumesnapshotResource).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		if snapshotUnstructured == nil {
			continue
		}
		for _, snapshotItem := range snapshotUnstructured.Items {
			snapshot := volumesnapshotv1.VolumeSnapshot{}

			runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotItem.Object, &snapshot)
			listSnapshots = append(listSnapshots, snapshot)
		}
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
