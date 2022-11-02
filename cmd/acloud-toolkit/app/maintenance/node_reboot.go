package maintenance

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/maintenance"
)

func newNodeRebootOptions() *maintenanceTaskOptions {
	return &maintenanceTaskOptions{}
}

func AddNodeRebootCreateFlags(flagSet *flag.FlagSet, opts *maintenanceTaskOptions) {
	flagSet.BoolVarP(&opts.allNodes, "all", "A", false, "reboot all nodes within the cluster")
	flagSet.DurationVar(&opts.timeout, "timeout", 0*time.Second, "The length of time to wait before giving up, zero means infinite")
}

// NewNodeRebootCmd returns the Cobra Bootstrap sub command
func NewNodeRebootCmd(upgradeOptions *maintenanceTaskOptions) *cobra.Command {
	if upgradeOptions == nil {
		upgradeOptions = newNodeRebootOptions()
	}

	var cmd = &cobra.Command{
		Use:   "node-reboot <node>",
		Args:  cobra.MinimumNArgs(0),
		Short: "reboot a kubernetes node within an Avisi Cloud Kubernetes cluster if required",
		Long: `reboot a kubernetes node within an Avisi Cloud Kubernetes cluster. This is only performed if the file '/var/run/reboot-required' is present on the host machine.
`,
		Example: `
acloud-toolkit maintenance node-reboot mynode
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			return runMaintenanceTask(ctx, args, *upgradeOptions, func(ctx context.Context, client *maintenance.NodeUpgraderClient, nodeName string) error {
				return client.RebootNode(ctx, nodeName)
			})
		},
	}

	AddNodeRebootCreateFlags(cmd.Flags(), upgradeOptions)

	return cmd
}
