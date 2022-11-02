package maintenance

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/maintenance"
)

func newmaintenanceUpgradeOptions() *maintenanceTaskOptions {
	return &maintenanceTaskOptions{}
}

func AddMaintenanceUpgradeCreateFlags(flagSet *flag.FlagSet, opts *maintenanceTaskOptions) {
	flagSet.BoolVarP(&opts.allNodes, "all", "A", false, "upgrade all nodes within the cluster")
	flagSet.DurationVar(&opts.timeout, "timeout", 0*time.Second, "The length of time to wait before giving up, zero means infinite")
}

// NewNodeUpgradeCmd returns the Cobra Bootstrap sub command
func NewNodeUpgradeCmd(upgradeOptions *maintenanceTaskOptions) *cobra.Command {
	if upgradeOptions == nil {
		upgradeOptions = newmaintenanceUpgradeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "node-upgrade <node>",
		Args:  cobra.MinimumNArgs(0),
		Short: "upgrade a kubernetes node within an Avisi Cloud Kubernetes cluster (Bring Your Own Node only)",
		Long: `Upgrade a kubernetes node within an Avisi Cloud Kubernetes cluster that is running with Bring Your Own Node enabled.

This command will upgrade both the Container Runtime and Kubelet version of a node.
`,
		Example: `
acloud-toolkit maintenance node-upgrade mynode
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			return runMaintenanceTask(ctx, args, *upgradeOptions, func(ctx context.Context, client *maintenance.NodeUpgraderClient, nodeName string) error {
				return client.UpgradeNode(ctx, nodeName)
			})
		},
	}

	AddMaintenanceUpgradeCreateFlags(cmd.Flags(), upgradeOptions)

	return cmd
}
