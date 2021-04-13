package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/ame/restoresnapshot"
)

type listOptions struct {
	sourceNamespace string
}

func newListOptions() *listOptions {
	return &listOptions{}
}

func AddListFlags(flagSet *flag.FlagSet, opts *listOptions) {
	flagSet.StringVar(&opts.sourceNamespace, "source-namespace", "default", "")
}

// NewListCmd returns the Cobra Bootstrap sub command
func NewListCmd(runOptions *listOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newListOptions()
	}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List a snapshot",
		Long:  `list a snapshot`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return restoresnapshot.List(runOptions.sourceNamespace)
		},
	}

	AddListFlags(cmd.Flags(), runOptions)

	return cmd
}
