package volumes

import (
	"context"
	"time"

	migrate_volume "gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/migrate-volume"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type batchMigrateVolumeOptions struct {
	sourceStorageClassName string
	targetStorageClassName string
	targetNamespace        string
	timeout                int32
	dryRun                 bool
	nodeSelector           []string
	migrationMode          string
	migrationFlags         string
}

func NewBatchMigrateVolumeOptions() *batchMigrateVolumeOptions {
	return &batchMigrateVolumeOptions{}
}

func AddBatchMigrateVolumeOptions(flagSet *flag.FlagSet, opts *batchMigrateVolumeOptions) {
	flagSet.StringVarP(&opts.sourceStorageClassName, "source-storage-class", "s", "", "name of the source storageclass")
	flagSet.StringVarP(&opts.targetStorageClassName, "target-storage-class", "t", "", "name of the target storageclass")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "default", "Namespace where the migrate job will be executed")
	flagSet.Int32Var(&opts.timeout, "timeout", 300, "Timeout of the context in minutes")
	flagSet.BoolVar(&opts.dryRun, "dry-run", false, "Perform a dry run of the batch migrate")
	flagSet.StringSliceVar(&opts.nodeSelector, "node-selector", []string{}, "comma separated list of node labels used for nodeSelector of the migration job")
	flagSet.StringVarP(&opts.migrationMode, "migration-mode", "m", "rsync", "Migration mode to use. Options: rsync, rclone")
	flagSet.StringVarP(&opts.migrationFlags, "migration-flags", "f", "", "Additional flags to pass to the migration tool")
}

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewBatchMigrateVolumeCmd(runOptions *batchMigrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = NewBatchMigrateVolumeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "batch-migrate",
		Short: "Batch migrate all volumes within a namespace to another storage class",
		Long:  `Batch migrate all volumes from a source storage class within a namespace to another storage class. For each PVC that has the source storage class within the namespace, this will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume(s) will remain available within the cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var cancel context.CancelFunc
			if runOptions.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			}
			defer cancel()
			return migrate_volume.BatchMigrateVolumes(ctx, migrate_volume.BatchMigrateOptions{
				SourceStorageClassName: runOptions.sourceStorageClassName,
				TargetStorageClassName: runOptions.targetStorageClassName,
				TargetNamespace:        runOptions.targetNamespace,
				Timeout:                runOptions.timeout,
				DryRun:                 runOptions.dryRun,
				MigrationMode:          migrate_volume.MigrationMode(runOptions.migrationMode),
			})
		},
	}

	AddBatchMigrateVolumeOptions(cmd.Flags(), runOptions)
	return cmd
}
