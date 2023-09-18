package snapshot

import (
	"context"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/snapshots"
)

type importOptions struct {
	namespace            string
	snapshotName         string
	snapshotStorageClass string
}

func newImportOptions() *importOptions {
	return &importOptions{}
}

func AddImportFlags(flagSet *flag.FlagSet, opts *importOptions) {
	flagSet.StringVar(&opts.namespace, "namespace", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.StringVar(&opts.snapshotName, "name", "", "name of the snapshot")
	flagSet.StringVar(&opts.snapshotStorageClass, "snapshot-storage-class", "", "")
}

func NewImportCmd(runOptions *importOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newImportOptions()
	}

	var cmd = &cobra.Command{
		Use:   "import <snapshot>",
		Args:  cobra.ExactArgs(1),
		Short: "Import raw Snapshot ID into a CSI snapshot.",
		Long: `This command creates Kubernetes CSI snapshot resources using a snapshot ID from the backend storage, for example AWS EBS, or Ceph RBD.
		`,
		Example: `
acloud-toolkit snapshot import --name example snap-12345
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return snapshots.ImportSnapshotFromRawID(context.Background(), runOptions.snapshotName, runOptions.namespace, runOptions.snapshotStorageClass, args[0])
		},
	}

	AddImportFlags(cmd.Flags(), runOptions)

	return cmd
}
