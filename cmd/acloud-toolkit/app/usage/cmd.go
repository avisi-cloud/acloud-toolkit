package usage

import (
	"github.com/spf13/cobra"
)

// NewUsageCmd returns cobra.Command to run the acloud-toolkit usage sub command
func NewUsageCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "usage",
		Short:   "usage for working with Kubernetes resource usage (experimental)",
		Long:    "usage for working with Kubernetes resource usages to automate various workflows (experimental)",
		Aliases: []string{"usages"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewListCmd(nil))
	return cmds
}
