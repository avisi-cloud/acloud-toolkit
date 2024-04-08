package snapshot

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/prune"
)

type pruneOptions struct {
	sourceNamespace string
	allNamespaces   bool
	dryRun          bool
	pruneImported   bool
	timeout         time.Duration
	minAge          time.Duration
}

func newPruneOptions() *pruneOptions {
	return &pruneOptions{}
}

func AddPruneFlags(flagSet *flag.FlagSet, opts *pruneOptions) {
	flagSet.StringVarP(&opts.sourceNamespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "Return results for all namespaces")
	flagSet.BoolVar(&opts.dryRun, "dry-run", true, "Perform a dry run of snapshot prune")
	flagSet.BoolVar(&opts.pruneImported, "prune-imported", false, "Prune snapshots imported from another cluster using acloud-toolkit, will ignore deletion policy")
	flagSet.DurationVar(&opts.minAge, "min-age", 0, "Minimum age of snapshots to be pruned")
	flagSet.DurationVarP(&opts.timeout, "timeout", "t", 10*time.Minute, "Duration to wait for the pruned snapshots to complete")
}

// NewPruneCmd returns the Cobra Bootstrap sub command
func NewPruneCmd(runOptions *pruneOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newPruneOptions()
	}

	var cmd = &cobra.Command{
		Use:   "prune",
		Short: "Prune removes any unused and released snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), runOptions.timeout)
			defer cancel()

			return prune.PruneSnapshots(ctxWithTimeout, runOptions.sourceNamespace, runOptions.allNamespaces, runOptions.pruneImported, runOptions.dryRun, runOptions.minAge)
		},
	}

	AddPruneFlags(cmd.Flags(), runOptions)

	return cmd
}
