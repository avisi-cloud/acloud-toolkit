package cmd

import (
	"context"
	"time"

	migrate_volume "gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/ame/migrate-volume"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type batchMigrateVolumeOptions struct {
	sourceStorageClassName string
	targetStorageClassName string
	targetNamespace        string
	timeout                int32
	dryRun                 bool
}

func NewBatchMigrateVolumeOptions() *batchMigrateVolumeOptions {
	return &batchMigrateVolumeOptions{}
}

func AddBatchMigrateVolumeOptions(flagSet *flag.FlagSet, opts *batchMigrateVolumeOptions) {
	flagSet.StringVarP(&opts.sourceStorageClassName, "source-storage-class", "s", "", "name of the source storageclass")
	flagSet.StringVarP(&opts.targetStorageClassName, "target-storage-class", "t", "", "name of the target storageclass")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "default", "Namespace where the migrate job will be executed")
	flagSet.Int32Var(&opts.timeout, "timeout", 60, "Timeout of the context in minutes")
	flagSet.BoolVar(&opts.dryRun, "dry-run", false, "Perform a dry run of the batch migrate")
}

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewBatchMigrateVolumeCmd(ctx context.Context, runOptions *batchMigrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = NewBatchMigrateVolumeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "batch-migrate",
		Short: "Batch migrate all volumes within a namespace to another storage class",
		Long:  `Batch migrate all volumes from a source storage class within a namespace to another storage class. For each PVC that has the source storage class within the namespace, this will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			defer cancel()
			return migrate_volume.BatchMigrateVolumes(ctxWithTimeout, runOptions.sourceStorageClassName, runOptions.targetStorageClassName, runOptions.targetNamespace, runOptions.dryRun)
		},
	}

	AddBatchMigrateVolumeOptions(cmd.Flags(), runOptions)

	return cmd
}
