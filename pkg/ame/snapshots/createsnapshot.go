package snapshots

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v8/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/utils/ptr"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
)

// SnapshotCreate creates a snapshot of a single PVC. The snapshot is labeled with SnapshotBatchLabelKey=batchID.
func SnapshotCreate(ctx context.Context, snapshotName string, targetNamespace string, targetName string, snapshotClassName string, batchID string) error {
	if err := ValidateSnapshotBatchID(batchID); err != nil {
		return err
	}
	if targetNamespace == "" {
		kubeconfig, err := k8s.GetClientConfig()
		if err != nil {
			return err
		}
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		targetNamespace = contextNamespace
	}

	submitted, err := createSnapshot(ctx, snapshotName, targetNamespace, targetName, snapshotClassName, batchID)
	if submitted {
		printSnapshotBatchDeleteHint(targetNamespace, batchID)
	}
	return err
}

// createSnapshot creates the VolumeSnapshot and waits until it is ready to use.
// submitted reports whether the VolumeSnapshot was created in the cluster, even if waiting for it failed afterwards.
func createSnapshot(ctx context.Context, snapshotName string, targetNamespace string, targetName string, snapshotClassName string, batchID string) (bool, error) {
	kubeconfig, err := k8s.GetClientConfig()
	if err != nil {
		return false, err
	}
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return false, err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return false, err
	}
	k8sclient, err := k8s.GetClientWithConfig(config)
	if err != nil {
		return false, err
	}
	if targetNamespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return false, err
		}
		targetNamespace = contextNamespace
	}

	// wait until PVC has a persistent volume
	pvc, err := k8sclient.CoreV1().PersistentVolumeClaims(targetNamespace).Get(ctx, targetName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	fmt.Printf("PVC with volume %s found...\n", pvc.Spec.VolumeName)

	var volumeSnapshotClassName *string
	if strings.TrimSpace(snapshotClassName) != "" {
		volumeSnapshotClassName = ptr.To(strings.TrimSpace(snapshotClassName))
	}
	// convert to the snapshot object
	volumeSnapshot := volumesnapshotv1.VolumeSnapshot{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "snapshot.storage.k8s.io/v1",
			Kind:       "VolumeSnapshot",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: targetNamespace,
			Labels: map[string]string{
				createdByLabelKey:     createdByLabelValue,
				SnapshotBatchLabelKey: batchID,
			},
		},
		Spec: volumesnapshotv1.VolumeSnapshotSpec{
			Source: volumesnapshotv1.VolumeSnapshotSource{
				PersistentVolumeClaimName: ptr.To(targetName),
			},
			VolumeSnapshotClassName: volumeSnapshotClassName,
		},
	}
	snapshotUnstructued, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&volumeSnapshot)
	if err != nil {
		return false, err
	}

	fmt.Printf("Creating snapshot %q for PVC %q...\n", snapshotName, targetName)
	_, err = client.Resource(volumesnapshotResource).Namespace(targetNamespace).Create(ctx, &unstructured.Unstructured{
		Object: snapshotUnstructued,
	}, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	fmt.Printf("Snapshot %q created for PVC %q...\n", snapshotName, targetName)
	fmt.Printf("Waiting until %q is ready for use...\n", snapshotName)
	for {
		snapshotUnstructued, err := client.Resource(volumesnapshotResource).Namespace(targetNamespace).Get(ctx, snapshotName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}

		// convert to the snapshot object
		snapshot := volumesnapshotv1.VolumeSnapshot{}
		runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotUnstructued.Object, &snapshot)

		// Check that the snapshot is ready to use
		if snapshot.Status != nil && snapshot.Status.ReadyToUse != nil && *snapshot.Status.ReadyToUse {
			fmt.Printf("Snapshot %q is ready for use...\n", snapshotName)
			break
		}
		select {
		case <-ctx.Done():
			return true, ctx.Err()
		case <-time.After(1 * time.Second):
			continue
		}
	}
	fmt.Printf("Snapshot %q completed...\n", snapshotName)

	return true, nil
}

// SnapshotCreateAllInNamespace creates a snapshot for every PVC in the namespace. All snapshots are labeled with
// SnapshotBatchLabelKey=batchID.
func SnapshotCreateAllInNamespace(ctx context.Context, targetNamespace, snapshotClassName, prefix string, concurrentLimit int, batchID string) error {
	if err := ValidateSnapshotBatchID(batchID); err != nil {
		return err
	}
	kubeconfig, err := k8s.GetClientConfig()
	if err != nil {
		return err
	}
	config, err := kubeconfig.ClientConfig()
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
	pvcs, err := k8sclient.CoreV1().PersistentVolumeClaims(targetNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var anySubmitted atomic.Bool
	errCh := make(chan error, len(pvcs.Items))
	createdSnapshots := make([]string, 0, len(pvcs.Items))

	// Create a worker pool with a limit on the number of concurrent snapshot creation operations
	workerPool := make(chan struct{}, concurrentLimit)

	suffix, _ := k8s.GenerateRandomString(5)

	for _, pvc := range pvcs.Items {
		wg.Add(1)
		workerPool <- struct{}{} // Reserve a spot in the worker pool

		go func(pvcName string) {
			defer wg.Done()
			defer func() { <-workerPool }() // Release the spot in the worker pool

			var snapshotName string
			if prefix != "" {
				snapshotName = k8s.FormatKubernetesNameCustomSuffix(fmt.Sprintf("%s-%s", prefix, pvcName), k8s.MaxKubernetesNameLength, suffix)
			} else {
				snapshotName = k8s.FormatKubernetesNameCustomSuffix(pvcName, k8s.MaxKubernetesNameLength, suffix)
			}

			submitted, err := createSnapshot(ctx, snapshotName, targetNamespace, pvcName, snapshotClassName, batchID)
			if submitted {
				anySubmitted.Store(true)
			}
			if err != nil {
				errCh <- err
			} else {
				mu.Lock()
				createdSnapshots = append(createdSnapshots, snapshotName)
				mu.Unlock()
			}
		}(pvc.Name)
	}

	wg.Wait()
	close(errCh)

	var finalErr error
	for err := range errCh {
		if err != nil {
			finalErr = err
			fmt.Printf("Error creating snapshot: %v\n", err)
		}
	}

	fmt.Println("\nCreated snapshots:")
	for _, snapshot := range createdSnapshots {
		fmt.Println(snapshot)
	}
	if anySubmitted.Load() {
		printSnapshotBatchDeleteHint(targetNamespace, batchID)
	}

	return finalErr
}
