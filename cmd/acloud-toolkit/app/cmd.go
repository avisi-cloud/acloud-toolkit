package app

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app/docs"
	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app/nodes"
	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app/snapshot"
	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app/version"
	"gitlab.avisi.cloud/ame/acloud-toolkit/cmd/acloud-toolkit/app/volumes"
)

// Execute runs the acloud-toolkit application
func Execute() error {
	cmd := NewACloudToolKitCmd(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}

// NewACloudToolKitCmd returns cobra.Command to run the acloud-toolkit command
func NewACloudToolKitCmd(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "acloud-toolkit",
		Short: "acloud-toolkit for working with Kubernetes",
		Long:  "acloud-toolkit for working with Kubernetes to automate various common tasks and workflows",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmds.ResetFlags()

	cmds.AddCommand(version.NewVersionCmd())
	cmds.AddCommand(snapshot.NewSnapshotCmd())
	cmds.AddCommand(volumes.NewStorageCmd())
	cmds.AddCommand(nodes.NewNodesCmd())
	cmds.AddCommand(docs.NewOpenDocs())

	return cmds
}
