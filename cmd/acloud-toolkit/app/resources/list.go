package resources

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
		Use:   "list",
		Short: "resources list displays resource requests and limits within the namespace agregrated by deployment or statefulset (experimental)",
		Long: `resources list displays resource requests and limits within the namespace agregrated by deployment or statefulset. This is experimental functionality. The displayed CPU and Memory are a sum of all resource requests and limits for each container within a pod.

For example, a pod with two containers of each 100Mi Memory limits will show 200Mi for Memory limits.`,
		Aliases: []string{"ls"},
		Example: `
# List resource limits within a specific namespace

❯ ./bin/acloud-toolkit resources list -n nginx-ingress
NAMESPACE               TYPE            NAME                            CONTAINERS      REPLICAS        MEMORY 
nginx-ingress           Deployment      ingress-nginx-controller        1               2               150Mi 

# List resource limits for all namespaces

❯ ./bin/acloud-toolkit resources list -A
...

		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return resourceusage.List(runOptions.sourceNamespace, runOptions.allNamespaces)
		},
	}

	AddListFlags(cmd.Flags(), runOptions)

	return cmd
}
