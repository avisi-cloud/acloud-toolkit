package snapshots

import (
	"context"
	"errors"
	"fmt"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
)

type DeleteByBatchOpts struct {
	BatchID       string
	Namespace     string
	AllNamespaces bool
	DryRun        bool
}

// DeleteByBatch deletes all VolumeSnapshots labeled with SnapshotBatchLabelKey=opts.BatchID.
func DeleteByBatch(ctx context.Context, opts DeleteByBatchOpts) error {
	if err := ValidateSnapshotBatchID(opts.BatchID); err != nil {
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
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	k8sclient, err := k8s.GetClientWithConfig(config)
	if err != nil {
		return err
	}

	namespaces, err := resolveNamespaces(ctx, kubeconfig, k8sclient, opts.Namespace, opts.AllNamespaces)
	if err != nil {
		return err
	}
	return deleteSnapshotsByBatch(ctx, client, namespaces, opts.BatchID, opts.DryRun)
}

func deleteSnapshotsByBatch(ctx context.Context, client dynamic.Interface, namespaces []string, batchID string, dryRun bool) error {
	selector := SnapshotBatchLabelSelector(batchID)
	snapshotsToDelete, err := listVolumeSnapshots(ctx, client, namespaces, selector)
	if err != nil {
		return err
	}
	if len(snapshotsToDelete) == 0 {
		fmt.Printf("No snapshots found with label %s\n", selector)
		return nil
	}

	if dryRun {
		fmt.Printf("Snapshots that would be deleted (dry-run), label %s:\n", selector)
		for _, snapshot := range snapshotsToDelete {
			fmt.Printf("  %s/%s\n", snapshot.GetNamespace(), snapshot.GetName())
		}
		return nil
	}

	fmt.Printf("Deleting %d snapshot(s) with label %s...\n", len(snapshotsToDelete), selector)
	var errs []error
	deleted := 0
	for _, snapshot := range snapshotsToDelete {
		err := client.Resource(volumesnapshotResource).Namespace(snapshot.GetNamespace()).Delete(ctx, snapshot.GetName(), metav1.DeleteOptions{
			Preconditions: metav1.NewUIDPreconditions(string(snapshot.GetUID())),
		})
		if err != nil && !kubeerrors.IsNotFound(err) {
			fmt.Printf("  failed to delete %s/%s: %v\n", snapshot.GetNamespace(), snapshot.GetName(), err)
			errs = append(errs, fmt.Errorf("failed to delete snapshot %s/%s: %w", snapshot.GetNamespace(), snapshot.GetName(), err))
			continue
		}
		deleted++
		fmt.Printf("  deleted %s/%s\n", snapshot.GetNamespace(), snapshot.GetName())
	}
	fmt.Printf("Deleted %d of %d snapshot(s)\n", deleted, len(snapshotsToDelete))

	return errors.Join(errs...)
}
