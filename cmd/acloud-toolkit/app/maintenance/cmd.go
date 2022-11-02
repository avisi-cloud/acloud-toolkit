package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/maintenance"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewMaintenanceCmd returns cobra.Command to run the acloud-toolkit Maintenance sub command
func NewMaintenanceCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "maintenance",
		Short:   "Perform maintenance actions on Kubernetes clusters",
		Long:    "Perform maintenance actions on Kubernetes clusters",
		Aliases: []string{"maintenances"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewDrainCmd(nil))
	cmds.AddCommand(NewNodeUpgradeCmd(nil))
	cmds.AddCommand(NewNodeRebootCmd(nil))
	return cmds
}

type maintenanceTaskOptions struct {
	allNodes bool
	timeout  time.Duration
}

func runMaintenanceTask(ctx context.Context, args []string, opts maintenanceTaskOptions, task func(context.Context, *maintenance.NodeUpgraderClient, string) error) error {
	k8sClient := k8s.GetClientOrDie()
	client := maintenance.NewNodeMaintenanceClient(k8sClient, maintenance.Opts{})
	if client == nil {
		return fmt.Errorf("failed to create maintenance client")
	}

	if len(args) > 0 {
		for _, nodeName := range args {
			if err := task(ctx, client, nodeName); err != nil {
				return err
			}
		}
		return nil
	}
	if opts.allNodes {
		nodeList, err := k8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list all nodes: %s", err.Error())
		}
		for _, node := range nodeList.Items {
			if err := task(ctx, client, node.GetName()); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("missing node name")
	}
	return nil
}
