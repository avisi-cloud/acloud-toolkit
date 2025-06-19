package volumes

import (
	"context"
	_ "embed"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/syncvolume"
)

type syncVolumeOptions struct {
	syncvolume.SyncVolumeJobOptions
	timeout int32
}

func newSyncVolumeOptions() *syncVolumeOptions {
	return &syncVolumeOptions{}
}

func AddSyncVolumeOptions(flagSet *flag.FlagSet, opts *syncVolumeOptions) {
	flagSet.Int32VarP(&opts.timeout, "timeout", "t", 60, "timeout of the context in minutes")
	flagSet.StringVar(&opts.SourcePVCName, "source-pvc", "", "name of the source persitentvolumeclaim")
	flagSet.StringVar(&opts.TargetPVCName, "target-pvc", "", "name of the target persitentvolumeclaim")
	flagSet.StringVarP(&opts.Namespace, "namespace", "n", "", "namespace where the sync job will be executed")
	flagSet.BoolVar(&opts.RetainJob, "retain-job", false, "retain the job after completion")
	flagSet.BoolVar(&opts.CreateNewPVC, "create-pvc", false, "create a new PVC if the target PVC does not exist")
	flagSet.StringVar(&opts.NewStorageClassName, "storageclass", "ebs-restore", "name of the storageclass to use for the new PVC")
	flagSet.Int32Var(&opts.TtlSecondsAfterFinished, "ttl", 3600, "time to live in seconds after the job has finished, requires --retain-job to be true")
	flagSet.Int64Var(&opts.NewSize, "new-size", syncvolume.USE_EQUAL_SIZE, "use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC")
}

//go:embed examples/sync.txt
var syncExamples string

// NewSyncVolumeCmd returns the Cobra Bootstrap sub command
func NewSyncVolumeCmd(runOptions *syncVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newSyncVolumeOptions()
	}

	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Sync a volume to another existing volume, or create a new volume",
		Long:    `Sync a volume to another existing volume, or create a new volume. This will create a new PVC using the target storage class or use an existing one, and copy all file contents over to the new volume using rsync. The existing persistent volume and persistent volume claim will remain available in the cluster.`,
		Example: syncExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), time.Duration(runOptions.timeout)*time.Minute)
			defer cancel()

			runOptions.ExtraRsyncArgs = args
			return syncvolume.SyncVolumeJob(ctxWithTimeout, runOptions.SyncVolumeJobOptions)
		},
	}

	AddSyncVolumeOptions(cmd.Flags(), runOptions)

	cmd.MarkFlagRequired("source-pvc")
	cmd.MarkFlagRequired("target-pvc")

	return cmd
}
