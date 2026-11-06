package snapshots

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
)

const (
	createdByLabelKey   = "app.kubernetes.io/created-by"
	createdByLabelValue = "acloud-toolkit"

	// SnapshotBatchLabelKey is set on every VolumeSnapshot created by a single `snapshot create` run,
	// so all snapshots of that run can be selected (and deleted) at once.
	SnapshotBatchLabelKey = "acloud-toolkit.k8s.avisi.cloud/snapshot-batch"
)

// NewSnapshotBatchID returns a unique value for the SnapshotBatchLabelKey label,
// formatted as <UTC timestamp>-<random suffix>, e.g. 20261106-142501-x7k2p.
func NewSnapshotBatchID() (string, error) {
	suffix, err := k8s.GenerateRandomString(k8s.DefaultSuffixLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate snapshot batch id: %w", err)
	}
	return time.Now().UTC().Format("20060102-150405") + "-" + suffix, nil
}

// ValidateSnapshotBatchID checks that batchID can be used as a Kubernetes label value.
func ValidateSnapshotBatchID(batchID string) error {
	if batchID == "" {
		return fmt.Errorf("snapshot batch id must not be empty")
	}
	if errs := validation.IsValidLabelValue(batchID); len(errs) > 0 {
		return fmt.Errorf("invalid snapshot batch id %q: %s", batchID, strings.Join(errs, "; "))
	}
	return nil
}

// SnapshotBatchLabelSelector returns the label selector matching all snapshots of the given batch.
func SnapshotBatchLabelSelector(batchID string) string {
	return SnapshotBatchLabelKey + "=" + batchID
}

func printSnapshotBatchDeleteHint(namespace, batchID string) {
	fmt.Printf("\nAll snapshots of this run are labeled %s\n", SnapshotBatchLabelSelector(batchID))
	fmt.Printf("Delete them in one go with:\n  acloud-toolkit snapshot delete --batch-id %s -n %s\n", batchID, namespace)
}
