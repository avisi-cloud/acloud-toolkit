package snapshot

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/snapshots"
)

type snapshotCreateOptions struct {
	persistentVolumeClaimName      string
	persistentVolumeClaimNamespace string
	snapshotCreateStorageClass     string
	timeout                        time.Duration
	allInNamespace                 bool
	prefix                         string
	concurrentLimit                int
	batchID                        string
}

func newSnapshotCreateOptions() *snapshotCreateOptions {
	return &snapshotCreateOptions{}
}

func AddSnapshotCreateFlags(flagSet *flag.FlagSet, opts *snapshotCreateOptions) {
	flagSet.StringVarP(&opts.persistentVolumeClaimNamespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.StringVarP(&opts.persistentVolumeClaimName, "pvc", "p", "", "Name of the PVC to snapshot. (required)")
	flagSet.StringVarP(&opts.snapshotCreateStorageClass, "snapshot-class", "s", "", "Name of the CSI volume snapshot class to use. Uses the default VolumeSnapshotClass by default")
	flagSet.DurationVarP(&opts.timeout, "timeout", "t", 60*time.Minute, "Duration to wait for the created snapshot to be ready for use")
	flagSet.BoolVar(&opts.allInNamespace, "all", false, "Create snapshots for all PVCs in the namespace, and use pvc name as snapshot name")
	flagSet.StringVar(&opts.prefix, "prefix", "", "Add a prefix seperated by `-` to the snapshot name when using --all")
	flagSet.IntVar(&opts.concurrentLimit, "concurrent-limit", 10, "Maximum number of concurrent snapshot creation operations")
	flagSet.StringVar(&opts.batchID, "batch-id", "", "Value for the "+snapshots.SnapshotBatchLabelKey+" label set on all created snapshots, to delete them later in one go. Generated (<timestamp>-<random>) by default; reuse a value to group multiple runs")
}

// NewSnapshotCreateCmd returns the Cobra Bootstrap sub command
func NewSnapshotCreateCmd(runOptions *snapshotCreateOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newSnapshotCreateOptions()
	}

	cmd := &cobra.Command{
		Use: "create <snapshot>",
		Args: func(cmd *cobra.Command, args []string) error {
			if !runOptions.allInNamespace && len(args) != 1 {
				return errors.New("requires 1 argument for snapshot name")
			}
			return nil
		},
		Short: "Create a snapshot of a Kubernetes PVC (persistent volume claim).",
		Long: `This command creates a snapshot of a Kubernetes PVC, allowing you to capture a point-in-time copy of the data stored in the PVC. Snapshots can be used for data backup, disaster recovery, and other purposes.

To create a snapshot, you need to provide the name of the PVC to snapshot, as well as a name for the snapshot. You can also specify a namespace if the PVC is not in the current namespace context. If no snapshot class is specified, the default snapshot class will be used.

All snapshots created by a single run are labeled with ` + "`" + snapshots.SnapshotBatchLabelKey + "=<batch-id>`" + `, so they can be deleted in one go afterwards. The batch id is generated unless --batch-id is set, and the matching delete command is printed when the run finishes.`,
		Example: `
# Create a snapshot of the PVC "my-pvc" with the name "my-snapshot":
acloud-toolkit snapshot create my-snapshot --pvc=my-pvc

#Create a snapshot of the PVC "my-pvc" with the name "my-snapshot" in the namespace "my-namespace":
acloud-toolkit snapshot create my-snapshot --pvc=my-pvc --namespace=my-namespace

# Create snapshots for all PVCs in the namespace "my-namespace":
acloud-toolkit snapshot create --all --namespace=my-namespace

# Create snapshots for all PVCs in the namespace "my-namespace" with a prefix "backup":
acloud-toolkit snapshot create --all --namespace=my-namespace --prefix=backup

# Create snapshots for all PVCs in the namespace "my-namespace" with a fixed batch id, then delete them all at once:
acloud-toolkit snapshot create --all --namespace=my-namespace --batch-id=before-upgrade
kubectl delete volumesnapshots -n my-namespace -l ` + snapshots.SnapshotBatchLabelKey + `=before-upgrade
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if runOptions.prefix != "" && !runOptions.allInNamespace {
				return errors.New("--prefix can only be set if --all is also provided")
			}
			if cmd.Flags().Changed("batch-id") {
				return snapshots.ValidateSnapshotBatchID(runOptions.batchID)
			}
			batchID, err := snapshots.NewSnapshotBatchID()
			if err != nil {
				return err
			}
			runOptions.batchID = batchID
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), runOptions.timeout)
			defer cancel()

			if runOptions.allInNamespace {
				return snapshots.SnapshotCreateAllInNamespace(ctxWithTimeout, runOptions.persistentVolumeClaimNamespace, runOptions.snapshotCreateStorageClass, runOptions.prefix, runOptions.concurrentLimit, runOptions.batchID)
			}

			return snapshots.SnapshotCreate(ctxWithTimeout, args[0], runOptions.persistentVolumeClaimNamespace, runOptions.persistentVolumeClaimName, runOptions.snapshotCreateStorageClass, runOptions.batchID)
		},
	}

	AddSnapshotCreateFlags(cmd.Flags(), runOptions)

	return cmd
}
