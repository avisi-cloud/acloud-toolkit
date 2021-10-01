package storage

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/restoresnapshot"
)

type snapshotCreateOptions struct {
	snapshotName               string
	sourceNamespace            string
	targetNamespace            string
	targetName                 string
	snapshotCreateStorageClass string
}

func newSnapshotCreateOptions() *snapshotCreateOptions {
	return &snapshotCreateOptions{}
}

func AddSnapshotCreateFlags(flagSet *flag.FlagSet, opts *snapshotCreateOptions) {
	flagSet.StringVar(&opts.snapshotName, "snapshot-name", "", "The name of the snapshot that will be created")
	flagSet.StringVarP(&opts.targetNamespace, "namespace", "n", "default", "Namespace of the PVC. Snapshot will be created within this namespace as well")
	flagSet.StringVarP(&opts.targetName, "pvc", "p", "", "Name of the persistent volume to snapshot")
	flagSet.StringVarP(&opts.snapshotCreateStorageClass, "snapshot-class", "s", "csi-aws-vsc", "CSI snapshot class")
}

// NewSnapshotCreateCmd returns the Cobra Bootstrap sub command
func NewSnapshotCreateCmd(runOptions *snapshotCreateOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newSnapshotCreateOptions()
	}

	var cmd = &cobra.Command{
		Use:   "create-snapshot",
		Short: "create-snapshot creates a snapshot for a pvc",
		Long:  `create-snapshot creates a snapshot for a pvc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return restoresnapshot.SnapshotCreate(runOptions.snapshotName, runOptions.targetName, runOptions.targetNamespace, runOptions.snapshotCreateStorageClass)
		},
	}

	AddSnapshotCreateFlags(cmd.Flags(), runOptions)

	return cmd
}
