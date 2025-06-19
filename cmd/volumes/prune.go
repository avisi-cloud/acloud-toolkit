package volumes

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/prune"
)

type volumePruneOptions struct {
	dryRun        bool
	allNamespaces bool
	pvcNamespace  string
}

func newvolumePruneOptions() *volumePruneOptions {
	return &volumePruneOptions{}
}

func AddvolumePruneFlags(flagSet *flag.FlagSet, opts *volumePruneOptions) {
	flagSet.BoolVar(&opts.dryRun, "dry-run", true, "Perform a dry run of volume prune")
	flagSet.BoolVarP(&opts.allNamespaces, "all", "A", false, "Prune volumes from all namespaces")
	flagSet.StringVarP(&opts.pvcNamespace, "namespace", "n", "", "Namespace to prune volumes from. Volume namespaces are cluster scoped, so the namespace is only used to filter the PVCs")
}

// NewvolumePruneCmd returns the Cobra Bootstrap sub command
func NewvolumePruneCmd(runOptions *volumePruneOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newvolumePruneOptions()
	}

	cmd := &cobra.Command{
		Use:   "prune <persistent-volume-claim>",
		Short: "Prune removes any unused and released persistent volumes",
		Long:  `The 'prune' command removes any released persistent volumes. By default it will run in dry-run mode, which will only show the volumes that would be pruned. Use the --dry-run=false flag to actually prune the volumes.`,
		Example: `
# See all persistent volumes that are set to Released that would be pruned
acloud-toolkit storage prune -A

# Prune all persistent volumes that are set to Released
acloud-toolkit storage prune -A --dry-run=false

# Prune all persistent volumes that are set to Released within a specific namespace
acloud-toolkit storage prune -n my-namespace --dry-run=false
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := prune.Volumes(cmd.Context(), prune.Opts{
				DryRun:        runOptions.dryRun,
				AllNamespaces: runOptions.allNamespaces,
				PvcNamespace:  runOptions.pvcNamespace,
			}); err != nil {
				return err
			}
			return nil
		},
	}

	AddvolumePruneFlags(cmd.Flags(), runOptions)

	return cmd
}
