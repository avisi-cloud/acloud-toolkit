package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCikCmd returns cobra.Command to run the cik command
func NewCikCmd(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "csi-snapshot-utils",
		Short: "csi-snapshot-utils for working with Kubernetes CSI snapshots",
		Long:  "csi-snapshot-utils for working with Kubernetes CSI snapshots to automate various workflows",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()

	cmds.AddCommand(NewVersionCmd())
	cmds.AddCommand(NewListCmd(nil))
	cmds.AddCommand(NewRestoreCmd(nil))

	return cmds
}
