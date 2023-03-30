package nodes

import (
	"github.com/spf13/cobra"
)

// NewNodesCmd returns cobra.Command to run the acloud-toolkit Maintenance sub command
func NewNodesCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "nodes",
		Short:   "Perform actions on Kubernetes cluster nodes",
		Long:    "Perform actions on Kubernetes cluster nodes",
		Aliases: []string{"node"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewDrainCmd(nil))
	return cmds
}
