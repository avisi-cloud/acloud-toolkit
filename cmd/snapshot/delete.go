package snapshot

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/snapshots"
)

type deleteOptions struct {
	batchID       string
	namespace     string
	allNamespaces bool
	dryRun        bool
}

func newDeleteOptions() *deleteOptions {
	return &deleteOptions{}
}

func AddDeleteFlags(flagSet *flag.FlagSet, opts *deleteOptions) {
	flagSet.StringVar(&opts.batchID, "batch-id", "", "Delete all snapshots of this batch (created by a single `snapshot create` run). (required)")
	flagSet.StringVarP(&opts.namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "Delete the snapshots of the batch in all namespaces")
	flagSet.BoolVar(&opts.dryRun, "dry-run", false, "Only show the snapshots that would be deleted")
}

// NewDeleteCmd returns the Cobra snapshot delete sub command
func NewDeleteCmd(runOptions *deleteOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newDeleteOptions()
	}

	cmd := &cobra.Command{
		Use:   "delete --batch-id <batch-id>",
		Short: "Delete all CSI snapshots created by a single snapshot create run",
		Long: `This command deletes all CSI snapshots of a batch. Every ` + "`acloud-toolkit snapshot create`" + ` run labels the snapshots it creates with ` + "`" + snapshots.SnapshotBatchLabelKey + "=<batch-id>`" + `, and prints the batch id when it finishes. The batch id of existing snapshots is shown in the Batch column of ` + "`acloud-toolkit snapshot list`" + `.

Whether the underlying snapshot data is removed as well depends on the deletion policy of the VolumeSnapshotClass.`,
		Example: `
# Show which snapshots of batch "20261106-142501-x7k2p" in the current namespace would be deleted:
acloud-toolkit snapshot delete --batch-id=20261106-142501-x7k2p --dry-run

# Delete all snapshots of batch "before-upgrade" in the namespace "my-namespace":
acloud-toolkit snapshot delete --batch-id=before-upgrade --namespace=my-namespace

# Delete all snapshots of batch "before-upgrade" in all namespaces:
acloud-toolkit snapshot delete --batch-id=before-upgrade -A
		`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return snapshots.DeleteByBatch(cmd.Context(), snapshots.DeleteByBatchOpts{
				BatchID:       runOptions.batchID,
				Namespace:     runOptions.namespace,
				AllNamespaces: runOptions.allNamespaces,
				DryRun:        runOptions.dryRun,
			})
		},
	}

	AddDeleteFlags(cmd.Flags(), runOptions)
	_ = cmd.MarkFlagRequired("batch-id")
	cmd.MarkFlagsMutuallyExclusive("namespace", "all-namespaces")

	return cmd
}
