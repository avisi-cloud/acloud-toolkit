package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"
)

// NewCSIUtilCmd returns cobra.Command to run the csi-utils command
func NewCSIUtilCmd(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "csi-utils",
		Short: "csi-utils for working with Kubernetes CSI",
		Long:  "csi-utils for working with Kubernetes CSI, volumes and snapshots to automate various workflows",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()

	cmds.AddCommand(NewVersionCmd())
	cmds.AddCommand(NewListCmd(nil))
	cmds.AddCommand(NewRestoreCmd(nil))
	cmds.AddCommand(NewSnapshotCreateCmd(nil))
	cmds.AddCommand(NewMigrateVolumeCmd(context.Background(), nil))
	return cmds
}
