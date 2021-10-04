package snapshot

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/restoresnapshot"
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
	flagSet.StringVarP(&opts.persistentVolumeClaimNamespace, "namespace", "n", "default", "Namespace of the PVC. Snapshot will be created within this namespace as well")
	flagSet.StringVarP(&opts.persistentVolumeClaimName, "pvc", "p", "", "Name of the persistent volume to snapshot")
	flagSet.StringVarP(&opts.snapshotCreateStorageClass, "snapshot-class", "s", "csi-aws-vsc", "CSI snapshot class")
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
			return restoresnapshot.SnapshotCreate(args[0], runOptions.persistentVolumeClaimNamespace, runOptions.persistentVolumeClaimName, runOptions.snapshotCreateStorageClass)
		},
	}

	AddSnapshotCreateFlags(cmd.Flags(), runOptions)

	return cmd
}
