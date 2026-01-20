package volumes

import (
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/prune"
)

type volumePruneOptions struct {
	dryRun              bool
	allNamespaces       bool
	pvcNamespace        string
	labelSelector       string
	minReleasedDuration time.Duration
	nameFilterPattern   string
}

func newVolumePruneOptions() *volumePruneOptions {
	return &volumePruneOptions{}
}

func AddVolumePruneFlags(flagSet *flag.FlagSet, opts *volumePruneOptions) {
	flagSet.BoolVar(&opts.dryRun, "dry-run", true, "Perform a dry run of volume prune")
	flagSet.BoolVarP(&opts.allNamespaces, "all", "A", false, "Prune volumes from all namespaces")
	flagSet.StringVarP(&opts.pvcNamespace, "namespace", "n", "", "Namespace to prune volumes from. Volume namespaces are cluster scoped, so the namespace is only used to filter the PVCs")
	flagSet.StringVarP(&opts.labelSelector, "label-selector", "l", "", "Label selector to filter the volumes to prune")
	flagSet.DurationVar(&opts.minReleasedDuration, "min-released-duration", 0, "Minimum duration since the volume was released")
	flagSet.StringVar(&opts.nameFilterPattern, "name-pattern", "", "Filter volumes by name using a regex pattern")
}

// NewVolumePruneCmd returns the Cobra Bootstrap sub command
func NewVolumePruneCmd(runOptions *volumePruneOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newVolumePruneOptions()
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

# Prune all persistent volumes that are set to Released within a specific namespace and name pattern
acloud-toolkit storage prune -n my-namespace --name-pattern "data-.*" --dry-run=false
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := prune.Volumes(cmd.Context(), prune.Opts{
				DryRun:              runOptions.dryRun,
				AllNamespaces:       runOptions.allNamespaces,
				PvcNamespace:        runOptions.pvcNamespace,
				LabelSelector:       runOptions.labelSelector,
				MinReleasedDuration: runOptions.minReleasedDuration,
				NameFilterPattern:   runOptions.nameFilterPattern,
			}); err != nil {
				return err
			}
			return nil
		},
	}

	AddVolumePruneFlags(cmd.Flags(), runOptions)

	return cmd
}
