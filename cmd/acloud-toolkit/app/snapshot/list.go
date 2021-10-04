package snapshot

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/restoresnapshot"
)

type listOptions struct {
	sourceNamespace string
	allNamespaces   bool
}

func newListOptions() *listOptions {
	return &listOptions{}
}

func AddListFlags(flagSet *flag.FlagSet, opts *listOptions) {
	flagSet.StringVarP(&opts.sourceNamespace, "namespace", "n", "", "return snapshots from a specific namespace. Default is the configured namespace in your kubecontext.")
	flagSet.BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "return results for all namespaces")
}

// NewListCmd returns the Cobra Bootstrap sub command
func NewListCmd(runOptions *listOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newListOptions()
	}

	var cmd = &cobra.Command{
		Use:     "list",
		Short:   "List CSI snapshots within the namespace",
		Long:    `List all available CSI snapshots within the namespace`,
		Aliases: []string{"ls"},
		Example: `
acloud-toolkit snapshot list
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return restoresnapshot.List(runOptions.sourceNamespace, runOptions.allNamespaces)
		},
	}

	AddListFlags(cmd.Flags(), runOptions)

	return cmd
}
