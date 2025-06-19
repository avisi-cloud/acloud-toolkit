package snapshot

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/snapshots"
)

type restoreOptions struct {
	sourceNamespace     string
	targetNamespace     string
	targetName          string
	restoreStorageClass string
	timeout             time.Duration
}

func newRestoreOptions() *restoreOptions {
	return &restoreOptions{}
}

func AddRestoreFlags(flagSet *flag.FlagSet, opts *restoreOptions) {
	flagSet.StringVar(&opts.sourceNamespace, "source-namespace", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
	flagSet.StringVar(&opts.targetNamespace, "restore-pvc-namespace", "", "")
	flagSet.StringVar(&opts.targetName, "restore-pvc-name", "", "")
	flagSet.StringVar(&opts.restoreStorageClass, "restore-storage-class", "", "")
	flagSet.DurationVarP(&opts.timeout, "timeout", "t", 10*time.Minute, "Duration to wait for the restored snapshot to complete")
}

// NewRestoreCmd returns the Cobra Bootstrap sub command
func NewRestoreCmd(runOptions *restoreOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newRestoreOptions()
	}

	cmd := &cobra.Command{
		Use:   "restore <snapshot>",
		Args:  cobra.ExactArgs(1),
		Short: "Restore a Kubernetes PVC from a CSI snapshot.",
		Long: `This command restores a Kubernetes PVC from a CSI snapshot. To restore a PVC, you need to provide the name of the snapshot, the name of the PVC to restore to, and the namespace of the target PVC. You can also specify a different namespace for the snapshot if needed.

By default, this command restores the PVC to the default storage class installed within the cluster. You can specify a different storage class if needed by using the --restore-storage-class option. Please note that this command requires the volume mode to be set to "Immediate".
		`,
		Example: `
acloud-toolkit snapshot restore my-snapshot --restore-pvc-name my-pvc --restore-storage-class ebs-restore
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), runOptions.timeout)
			defer cancel()
			return snapshots.RestoreSnapshot(ctxWithTimeout, args[0], runOptions.sourceNamespace, runOptions.targetName, runOptions.targetNamespace, runOptions.restoreStorageClass)
		},
	}

	AddRestoreFlags(cmd.Flags(), runOptions)

	return cmd
}
