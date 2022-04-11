package maintenance

import (
	"github.com/spf13/cobra"
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
	return cmds
}
