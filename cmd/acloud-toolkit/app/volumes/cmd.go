package volumes

import (
	"github.com/spf13/cobra"
)

// NewStorageCmd returns cobra.Command to run the acloud-toolkit storage sub command
func NewStorageCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "volumes",
		Aliases: []string{"storage"},
		Short:   "Various commands for working with Kubernetes CSI volumes",
		Long:    "Various commands for working with Kubernetes CSI volumes",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewMigrateVolumeCmd(nil))
	cmds.AddCommand(NewBatchMigrateVolumeCmd(nil))
	cmds.AddCommand(NewvolumeResizeCmd(nil))
	cmds.AddCommand(NewListCmd(nil))
	cmds.AddCommand(NewvolumePruneCmd(nil))
	cmds.AddCommand(NewSyncVolumeCmd(nil))
	return cmds
}
