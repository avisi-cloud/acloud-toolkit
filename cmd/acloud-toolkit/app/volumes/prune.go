package volumes

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/prune"
)

type volumePruneOptions struct {
	dryRun bool
}

func newvolumePruneOptions() *volumePruneOptions {
	return &volumePruneOptions{}
}

func AddvolumePruneFlags(flagSet *flag.FlagSet, opts *volumePruneOptions) {
	flagSet.BoolVar(&opts.dryRun, "dry-run", true, "Perform a dry run of volume prune")

}

// NewvolumePruneCmd returns the Cobra Bootstrap sub command
func NewvolumePruneCmd(runOptions *volumePruneOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newvolumePruneOptions()
	}

	var cmd = &cobra.Command{
		Use:   "prune <persistent-volume-claim>",
		Short: "Prune removes any unused and released persistent volumes",
		Long:  `The 'prune' command removes any released persistent volumes.`,
		Example: `
# Prune all persistent volumes that are set to Released
acloud-toolkit storage prune
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := prune.PruneVolumes(cmd.Context(), runOptions.dryRun); err != nil {
				return err
			}
			return nil
		},
	}

	AddvolumePruneFlags(cmd.Flags(), runOptions)

	return cmd
}
