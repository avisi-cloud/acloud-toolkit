package volumes

import (
	"context"
	_ "embed"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/migratevolume"
)

type migrateVolumeOptions struct {
	storageClassName string
	pvcName          string
	targetNamespace  string
	timeout          int32
	newSize          int64
	nodeSelector     []string
	migrationMode    string
	migrationFlags   string

	rsyncImage  string
	rcloneImage string
}

func newMigrateVolumeOptions() *migrateVolumeOptions {
	return &migrateVolumeOptions{}
}

func AddMigrateVolumeOptions(flagSet *flag.FlagSet, opts *migrateVolumeOptions) {
	flagSet.StringVarP(&opts.storageClassName, "storageClass", "s", "", "name of the new storageclass")
	flagSet.StringVarP(&opts.pvcName, "pvc", "p", "", "name of the persitentvolumeclaim")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "", "Namespace where the volume migrate job will be executed")
	flagSet.Int32VarP(&opts.timeout, "timeout", "t", 300, "Timeout of the context in minutes")
	flagSet.Int64Var(&opts.newSize, "new-size", migratevolume.UseEqualSize, "Use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC")
	flagSet.StringSliceVar(&opts.nodeSelector, "node-selector", []string{}, "comma separated list of node labels used for nodeSelector of the migration job")
	flagSet.StringVarP(&opts.migrationMode, "migration-mode", "m", "rsync", "Migration mode to use. Options: rsync, rclone. Default is rsync with rclone being newly introduced")
	flagSet.StringVarP(&opts.migrationFlags, "migration-flags", "f", "", "Additional flags to pass to the migration tool")

	// images
	flagSet.StringVar(&opts.rsyncImage, "rsync-image", migratevolume.DefaultRSyncContainerImage, "Image used for the rsync migration tool")
	flagSet.StringVar(&opts.rcloneImage, "rclone-image", migratevolume.DefaultRCloneContainerImage, "Image used for the rclone migration tool")
}

//go:embed examples/migrate.txt
var migrateExamples string

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewMigrateVolumeCmd(runOptions *migrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newMigrateVolumeOptions()
	}

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the filesystem on a persistent volume to another storage class",
		Long: `Migrate the filesystem on a persistent volume to another storage class.
This will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.

Migrate supports both rclone and rsync migration modes. The default mode is rsync.
- When using rsync, by default it uses the --archive flag. It will preserve all file permissions, timestamps, and ownerships.
- When using rclone a copy command is used. Use --metadata flag to preserve metadata.

It is recommended to utilize the migration-flag option to pass additional flags to the migration tool, such as --checksum or others and optmize the migration job for your specific use case.
`,
		Example: migrateExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var cancel context.CancelFunc
			if runOptions.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			}
			defer cancel()
			return migratevolume.StartMigrateVolumeJob(ctx, migratevolume.MigrationOptions{
				StorageClassName: runOptions.storageClassName,
				PVCName:          runOptions.pvcName,
				TargetNamespace:  runOptions.targetNamespace,
				NewSize:          runOptions.newSize,
				NodeSelector:     runOptions.nodeSelector,
				MigrationMode:    migratevolume.MigrationMode(runOptions.migrationMode),
				MigrationFlags:   runOptions.migrationFlags,

				RCloneImage: runOptions.rcloneImage,
				RyncImage:   runOptions.rsyncImage,
			})
		},
	}
	AddMigrateVolumeOptions(cmd.Flags(), runOptions)
	return cmd
}
