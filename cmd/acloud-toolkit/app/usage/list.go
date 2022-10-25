package usage

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/resourceusage"
)

type listOptions struct {
	sourceNamespace string
	allNamespaces   bool
}

func newListOptions() *listOptions {
	return &listOptions{}
}

func AddListFlags(flagSet *flag.FlagSet, opts *listOptions) {
	flagSet.StringVarP(&opts.sourceNamespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.BoolVarP(&opts.allNamespaces, "all-namespaces", "A", false, "return results for all namespaces")
}

// NewListCmd returns the Cobra Bootstrap sub command
func NewListCmd(runOptions *listOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newListOptions()
	}

	var cmd = &cobra.Command{
		Use:     "list",
		Short:   "List usages within the namespace (experimental)",
		Long:    `List usages within the namespace (experimental)`,
		Aliases: []string{"ls"},
		Example: `
acloud-toolkit usages list
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return resourceusage.List(runOptions.sourceNamespace, runOptions.allNamespaces)
		},
	}

	AddListFlags(cmd.Flags(), runOptions)

	return cmd
}
