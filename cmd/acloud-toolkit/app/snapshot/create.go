package snapshot

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/snapshots"
)

type snapshotCreateOptions struct {
	persistentVolumeClaimName      string
	persistentVolumeClaimNamespace string
	snapshotCreateStorageClass     string
}

func newSnapshotCreateOptions() *snapshotCreateOptions {
	return &snapshotCreateOptions{}
}

func AddSnapshotCreateFlags(flagSet *flag.FlagSet, opts *snapshotCreateOptions) {
	flagSet.StringVarP(&opts.persistentVolumeClaimNamespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.StringVarP(&opts.persistentVolumeClaimName, "pvc", "p", "", "Name of the persistent volume to snapshot")
	flagSet.StringVarP(&opts.snapshotCreateStorageClass, "snapshot-class", "s", "", "CSI volume snapshot class. If empty, use deafult volume snapshot class")
}

// NewSnapshotCreateCmd returns the Cobra Bootstrap sub command
func NewSnapshotCreateCmd(runOptions *snapshotCreateOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newSnapshotCreateOptions()
	}

	var cmd = &cobra.Command{
		Use:   "create <snapshot>",
		Args:  cobra.ExactArgs(1),
		Short: "create creates a snapshot for a pvc",
		Long:  `create creates a snapshot for a pvc`,
		Example: `
acloud-toolkit snapshot create my-snapshot --pvc my-pvc
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return snapshots.SnapshotCreate(args[0], runOptions.persistentVolumeClaimNamespace, runOptions.persistentVolumeClaimName, runOptions.snapshotCreateStorageClass)
		},
	}

	AddSnapshotCreateFlags(cmd.Flags(), runOptions)

	return cmd
}
