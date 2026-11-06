package snapshot

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/snapshots"
)

type listOptions struct {
	sourceNamespace      string
	allNamespaces        bool
	fetchSnapshotHandles bool
	batchID              string
}

func newListOptions() *listOptions {
	return &listOptions{}
}

func AddListFlags(flagSet *flag.FlagSet, opts *listOptions) {
	flagSet.StringVarP(&opts.sourceNamespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "return results for all namespaces")
	flagSet.BoolVarP(&opts.fetchSnapshotHandles, "handles", "S", true, "show snapshot content handle")
	flagSet.StringVar(&opts.batchID, "batch-id", "", "Only list snapshots of this batch (created by a single `snapshot create` run)")
}

// NewListCmd returns the Cobra Bootstrap sub command
func NewListCmd(runOptions *listOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newListOptions()
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available CSI snapshots within the current namespace",
		Long: `This command lists all available CSI snapshots within the current namespace. CSI snapshots are used to capture a point-in-time copy of a Kubernetes PVC, allowing you to preserve the data stored in the PVC for backup, disaster recovery, or other purposes.

By default, this command lists all snapshots in the current namespace. You can also specify a different namespace if needed.

The Batch column shows the batch id of snapshots created with ` + "`acloud-toolkit snapshot create`" + `. Use --batch-id to only list the snapshots of one batch.`,
		Aliases: []string{"ls"},
		Example: `
# List all available CSI snapshots within the current namespace:
acloud-toolkit snapshot list

# List all available CSI snapshots within the "my-namespace" namespace:
acloud-toolkit snapshot list --namespace=my-namespace

# List all available CSI snapshots within all namespaces:
acloud-toolkit snapshot list -A

# List the snapshots created by a single "snapshot create" run, in all namespaces:
acloud-toolkit snapshot list -A --batch-id=20261106-142501-x7k2p
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return snapshots.List(cmd.Context(), runOptions.sourceNamespace, runOptions.allNamespaces, runOptions.fetchSnapshotHandles, runOptions.batchID)
		},
	}

	AddListFlags(cmd.Flags(), runOptions)

	return cmd
}
