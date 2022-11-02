package resources

import (
	"github.com/spf13/cobra"
)

// NewResourcesCmd returns cobra.Command to run the acloud-toolkit usage sub command
func NewResourcesCmd() *cobra.Command {
	cmds := &cobra.Command{
		Use:     "resources",
		Short:   "Gather insight into resource usage and limits within a namespace (experimental)",
		Long:    `Gather insight into resource usage and limits within a namespace. This is experimental functionality.`,
		Aliases: []string{"resources"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewListCmd(nil))
	return cmds
}
