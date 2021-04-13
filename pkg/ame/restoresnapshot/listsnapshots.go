package restoresnapshot

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/k8s"

	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func List(sourceNamespace string) error {
	config, err := k8s.GetKubeConfigOrInCluster()
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	volumesnapshotRes := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshots"}
	snapshotUnstructued, err := client.Resource(volumesnapshotRes).Namespace(sourceNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, snapshot := range snapshotUnstructued.Items {
		fmt.Printf("%s\n", snapshot.GetName())
	}
	return nil
}
