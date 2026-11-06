package snapshots

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func newTestSnapshot(namespace, name, batchID string) *unstructured.Unstructured {
	snapshot := &unstructured.Unstructured{}
	snapshot.SetAPIVersion("snapshot.storage.k8s.io/v1")
	snapshot.SetKind("VolumeSnapshot")
	snapshot.SetNamespace(namespace)
	snapshot.SetName(name)
	if batchID != "" {
		snapshot.SetLabels(map[string]string{SnapshotBatchLabelKey: batchID})
	}
	return snapshot
}

func newFakeSnapshotClient(objects ...runtime.Object) *dynamicfake.FakeDynamicClient {
	return dynamicfake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(),
		map[schema.GroupVersionResource]string{volumesnapshotResource: "VolumeSnapshotList"},
		objects...,
	)
}

func remainingSnapshots(t *testing.T, client *dynamicfake.FakeDynamicClient, namespaces ...string) []string {
	t.Helper()
	names := []string{}
	for _, namespace := range namespaces {
		list, err := client.Resource(volumesnapshotResource).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
		require.NoError(t, err)
		for _, item := range list.Items {
			names = append(names, item.GetNamespace()+"/"+item.GetName())
		}
	}
	return names
}

func testSnapshots() []runtime.Object {
	return []runtime.Object{
		newTestSnapshot("ns-a", "a-batch1", "batch1"),
		newTestSnapshot("ns-a", "a-batch2", "batch2"),
		newTestSnapshot("ns-a", "a-unlabeled", ""),
		newTestSnapshot("ns-b", "b-batch1", "batch1"),
	}
}

func TestListVolumeSnapshotsByBatch(t *testing.T) {
	client := newFakeSnapshotClient(testSnapshots()...)

	snapshots, err := listVolumeSnapshots(context.Background(), client, []string{"ns-a", "ns-b"}, SnapshotBatchLabelSelector("batch1"))
	require.NoError(t, err)

	names := []string{}
	for _, snapshot := range snapshots {
		names = append(names, snapshot.GetNamespace()+"/"+snapshot.GetName())
	}
	assert.ElementsMatch(t, []string{"ns-a/a-batch1", "ns-b/b-batch1"}, names)
}

func TestDeleteSnapshotsByBatch(t *testing.T) {
	client := newFakeSnapshotClient(testSnapshots()...)

	require.NoError(t, deleteSnapshotsByBatch(context.Background(), client, []string{"ns-a"}, "batch1", false))

	assert.ElementsMatch(t,
		[]string{"ns-a/a-batch2", "ns-a/a-unlabeled", "ns-b/b-batch1"},
		remainingSnapshots(t, client, "ns-a", "ns-b"),
		"only batch1 snapshots in the selected namespace are deleted")
}

func TestDeleteSnapshotsByBatchDryRun(t *testing.T) {
	client := newFakeSnapshotClient(testSnapshots()...)

	require.NoError(t, deleteSnapshotsByBatch(context.Background(), client, []string{"ns-a", "ns-b"}, "batch1", true))

	assert.Len(t, remainingSnapshots(t, client, "ns-a", "ns-b"), 4, "dry-run deletes nothing")
}

func TestDeleteSnapshotsByBatchNoMatches(t *testing.T) {
	client := newFakeSnapshotClient(testSnapshots()...)

	require.NoError(t, deleteSnapshotsByBatch(context.Background(), client, []string{"ns-a", "ns-b"}, "unknown", false))

	assert.Len(t, remainingSnapshots(t, client, "ns-a", "ns-b"), 4)
}
