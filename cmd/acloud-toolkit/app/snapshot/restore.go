package snapshot

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/restoresnapshot"
)

type restoreOptions struct {
	sourceNamespace     string
	targetNamespace     string
	targetName          string
	restoreStorageClass string
}

func newRestoreOptions() *restoreOptions {
	return &restoreOptions{}
}

func AddRestoreFlags(flagSet *flag.FlagSet, opts *restoreOptions) {
	flagSet.StringVar(&opts.sourceNamespace, "source-namespace", "default", "")
	flagSet.StringVar(&opts.targetNamespace, "target-namespace", "default", "")
	flagSet.StringVar(&opts.targetName, "target-name", "", "")
	flagSet.StringVar(&opts.restoreStorageClass, "restore-storage-class", "ebs-restore", "")
}

// NewRestoreCmd returns the Cobra Bootstrap sub command
func NewRestoreCmd(runOptions *restoreOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newRestoreOptions()
	}

	var cmd = &cobra.Command{
		Use:   "restore <snapshot>",
		Args:  cobra.ExactArgs(1),
		Short: "Restore a snapshot",
		Long:  `restore a snapshot`,
		Example: `
acloud-toolkit snapshot restore my-snapshot --target-name my-pvc --restore-storage-class ebs-restore
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return restoresnapshot.Restore(args[0], runOptions.sourceNamespace, runOptions.targetName, runOptions.targetNamespace, runOptions.restoreStorageClass)
		},
	}

	AddRestoreFlags(cmd.Flags(), runOptions)

	return cmd
}
