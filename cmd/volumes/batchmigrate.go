package volumes

import (
	"context"
	_ "embed"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/migratevolume"
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
	preserveMetadata       bool

	rsyncImage  string
	rcloneImage string
}

func newBatchMigrateVolumeOptions() *batchMigrateVolumeOptions {
	return &batchMigrateVolumeOptions{}
}

func AddBatchMigrateVolumeOptions(flagSet *flag.FlagSet, opts *batchMigrateVolumeOptions) {
	flagSet.StringVarP(&opts.sourceStorageClassName, "source-storage-class", "s", "", "name of the source storageclass")
	flagSet.StringVarP(&opts.targetStorageClassName, "target-storage-class", "t", "", "name of the target storageclass")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "", "Namespace where the migrate job will be executed")
	flagSet.Int32Var(&opts.timeout, "timeout", 300, "Timeout of the context in minutes")
	flagSet.BoolVar(&opts.dryRun, "dry-run", false, "Perform a dry run of the batch migrate")
	flagSet.StringSliceVar(&opts.nodeSelector, "node-selector", []string{}, "comma separated list of node labels used for nodeSelector of the migration job")
	flagSet.StringVarP(&opts.migrationMode, "migration-mode", "m", "rsync", "Migration mode to use. Options: rsync, rclone")
	flagSet.StringVarP(&opts.migrationFlags, "migration-flags", "f", "", "Additional flags to pass to the migration tool")
	flagSet.BoolVar(&opts.preserveMetadata, "preserve-metadata", false, "Preserve the original metadata of the PVC")
	// images
	flagSet.StringVar(&opts.rsyncImage, "rsync-image", migratevolume.DefaultRSyncContainerImage, "Image used for the rsync migration tool")
	flagSet.StringVar(&opts.rcloneImage, "rclone-image", migratevolume.DefaultRCloneContainerImage, "Image used for the rclone migration tool")
}

//go:embed examples/batch-migrate.txt
var batchmigrateExamples string

// NewBatchMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewBatchMigrateVolumeCmd(runOptions *batchMigrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newBatchMigrateVolumeOptions()
	}

	cmd := &cobra.Command{
		Use:   "batch-migrate",
		Short: "Batch migrate all volumes within a namespace to another storage class",
		Long: `Batch migrate all volumes from a source storage class within a namespace to another storage class.
For each PVC that has the source storage class within the namespace, this will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume(s) will remain available within the cluster.

Batch migrate supports both rclone and rsync migration modes. The default mode is rsync.
- When using rsync, by default it uses the --archive flag. It will preserve all file permissions, timestamps, and ownerships.
- When using rclone a copy command is used. Use --metadata flag to preserve metadata.

It is recommended to utilize the migration-flag option to pass additional flags to the migration tool, such as --checksum or others and optmize the migration job for your specific use case.
		`,
		Example: batchmigrateExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var cancel context.CancelFunc
			if runOptions.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			}
			defer cancel()
			return migratevolume.BatchMigrateVolumes(ctx, migratevolume.BatchMigrateOptions{
				SourceStorageClassName: runOptions.sourceStorageClassName,
				TargetStorageClassName: runOptions.targetStorageClassName,
				TargetNamespace:        runOptions.targetNamespace,
				Timeout:                runOptions.timeout,
				DryRun:                 runOptions.dryRun,
				MigrationMode:          migratevolume.MigrationMode(runOptions.migrationMode),
				MigrationFlags:         runOptions.migrationFlags,
				NodeSelector:           runOptions.nodeSelector,
				PreserveMetadata:       runOptions.preserveMetadata,
				RSyncImage:             runOptions.rsyncImage,
				RCloneImage:            runOptions.rcloneImage,
			})
		},
	}

	AddBatchMigrateVolumeOptions(cmd.Flags(), runOptions)
	return cmd
}
