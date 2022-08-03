package snapshots

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/table"

	"k8s.io/client-go/dynamic"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func List(sourceNamespace string, allNamespaces bool) error {
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

	namespace, _, err := kubeconfig.Namespace()
	if err != nil {
		return err
	}
	if namespace == "" {
		return nil
	}
	if sourceNamespace != "" {
		namespace = sourceNamespace
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
		volumesnapshotRes := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"}
		snapshotUnstructured, err := client.Resource(volumesnapshotRes).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
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
			sourceName = *snapshot.Spec.Source.PersistentVolumeClaimName
		}
		contentName := ""
		if snapshot.Spec.Source.VolumeSnapshotContentName != nil {
			contentName = *snapshot.Spec.Source.VolumeSnapshotContentName
		}

		snapshotClassName := ""
		if snapshot.Spec.VolumeSnapshotClassName != nil {
			snapshotClassName = *snapshot.Spec.VolumeSnapshotClassName
		}

		body = append(body, []string{
			snapshot.GetNamespace(),
			snapshot.Name,
			sourceName,
			contentName,
			snapshot.Status.RestoreSize.String(),
			fmt.Sprintf("%v", snapshot.Status != nil && snapshot.Status.ReadyToUse != nil && *snapshot.Status.ReadyToUse),
			snapshotClassName,
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
	}, body)

	return nil
}
