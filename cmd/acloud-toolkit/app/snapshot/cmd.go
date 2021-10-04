package snapshot

import (
	"github.com/spf13/cobra"
)

// NewSnapshotCmd returns cobra.Command to run the acloud-toolkit Snapshot sub command
func NewSnapshotCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "snapshot",
		Short:   "snapshot for working with Kubernetes CSI snapshot",
		Long:    "snapshot for working with Kubernetes CSI snapshots to automate various workflows",
		Aliases: []string{"snapshots"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewListCmd(nil))
	cmds.AddCommand(NewRestoreCmd(nil))
	cmds.AddCommand(NewSnapshotCreateCmd(nil))
	return cmds
}
