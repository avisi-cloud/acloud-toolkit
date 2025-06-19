package snapshot

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/snapshots"
)

type snapshotCreateOptions struct {
	persistentVolumeClaimName      string
	persistentVolumeClaimNamespace string
	snapshotCreateStorageClass     string
	timeout                        time.Duration
	allInNamespace                 bool
	prefix                         string
	concurrentLimit                int
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

To create a snapshot, you need to provide the name of the PVC to snapshot, as well as a name for the snapshot. You can also specify a namespace if the PVC is not in the current namespace context. If no snapshot class is specified, the default snapshot class will be used.`,
		Example: `
# Create a snapshot of the PVC "my-pvc" with the name "my-snapshot":
acloud-toolkit snapshot create my-snapshot --pvc=my-pvc

#Create a snapshot of the PVC "my-pvc" with the name "my-snapshot" in the namespace "my-namespace":
acloud-toolkit snapshot create my-snapshot --pvc=my-pvc --namespace=my-namespace

# Create snapshots for all PVCs in the namespace "my-namespace":
acloud-toolkit snapshot create --all --namespace=my-namespace

# Create snapshots for all PVCs in the namespace "my-namespace" with a prefix "backup":
acloud-toolkit snapshot create --all --namespace=my-namespace --prefix=backup
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if runOptions.prefix != "" && !runOptions.allInNamespace {
				return errors.New("--prefix can only be set if --all is also provided")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), runOptions.timeout)
			defer cancel()

			if runOptions.allInNamespace {
				return snapshots.SnapshotCreateAllInNamespace(ctxWithTimeout, runOptions.persistentVolumeClaimNamespace, runOptions.snapshotCreateStorageClass, runOptions.prefix, runOptions.concurrentLimit)
			}

			return snapshots.SnapshotCreate(ctxWithTimeout, args[0], runOptions.persistentVolumeClaimNamespace, runOptions.persistentVolumeClaimName, runOptions.snapshotCreateStorageClass)
		},
	}

	AddSnapshotCreateFlags(cmd.Flags(), runOptions)

	return cmd
}
