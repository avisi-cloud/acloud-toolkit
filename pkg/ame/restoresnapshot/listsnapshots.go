package restoresnapshot

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/table"

	"k8s.io/client-go/dynamic"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func List(sourceNamespace string) error {
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

	volumesnapshotRes := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"}
	snapshotUnstructued, err := client.Resource(volumesnapshotRes).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	body := make([][]string, 0, len(snapshotUnstructued.Items))
	for _, snapshotItem := range snapshotUnstructued.Items {
		snapshot := volumesnapshotv1.VolumeSnapshot{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotItem.Object, &snapshot)

		sourceName := ""
		if snapshot.Spec.Source.PersistentVolumeClaimName != nil {
			sourceName = *snapshot.Spec.Source.PersistentVolumeClaimName
		}
		contentName := ""
		if snapshot.Spec.Source.VolumeSnapshotContentName != nil {
			contentName = *snapshot.Spec.Source.VolumeSnapshotContentName
		}

		body = append(body, []string{
			snapshot.Name,
			sourceName,
			contentName,
			snapshot.Status.RestoreSize.String(),
			fmt.Sprintf("%v", snapshot.Status != nil && snapshot.Status.ReadyToUse != nil && *snapshot.Status.ReadyToUse),
			*snapshot.Spec.VolumeSnapshotClassName,
		})
	}
	table.Print([]string{
		"Name",
		"Source",
		"Content",
		"Size",
		"Ready",
		"Classname",
	}, body)

	return nil
}
