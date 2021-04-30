package restoresnapshot

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/k8s"

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
	for _, snapshotItem := range snapshotUnstructued.Items {
		snapshot := volumesnapshotv1.VolumeSnapshot{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotItem.Object, &snapshot)
		fmt.Printf("%s\t\n", snapshot.GetName())
	}
	return nil
}
