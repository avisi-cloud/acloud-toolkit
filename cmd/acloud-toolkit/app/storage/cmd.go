package storage

import (
	"context"

	"github.com/spf13/cobra"
)

// NewStorageCmd returns cobra.Command to run the acloud-toolkit storage sub command
func NewStorageCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "storage",
		Short: "storage for working with Kubernetes CSI",
		Long:  "storage for working with Kubernetes CSI, volumes and snapshots to automate various workflows",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewMigrateVolumeCmd(context.Background(), nil))
	cmds.AddCommand(NewBatchMigrateVolumeCmd(context.Background(), nil))
	cmds.AddCommand(NewvolumeResizeCmd(nil))
	cmds.AddCommand(NewvolumePruneCmd(nil))
	return cmds
}
