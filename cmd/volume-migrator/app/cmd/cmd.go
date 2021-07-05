package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCikCmd returns cobra.Command to run the cik command
func NewCikCmd(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "volume-migrator",
		Short: "volume-migrator for moving Kubernetes volumes",
		Long:  "volume-migrator for moving Kubernetes volumes to automate various workflows",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
    cmds.AddCommand(NewMigrateVolumeCmd(nil))
	return cmds
}
