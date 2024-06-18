package volumes

import (
	"context"
	"time"

	migrate_volume "gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/migrate-volume"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
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
}

func NewMigrateVolumeOptions() *migrateVolumeOptions {
	return &migrateVolumeOptions{}
}

func AddMigrateVolumeOptions(flagSet *flag.FlagSet, opts *migrateVolumeOptions) {
	flagSet.StringVarP(&opts.storageClassName, "storageClass", "s", "", "name of the new storageclass")
	flagSet.StringVarP(&opts.pvcName, "pvc", "p", "", "name of the persitentvolumeclaim")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "default", "Namespace where the volume migrate job will be executed")
	flagSet.Int32VarP(&opts.timeout, "timeout", "t", 300, "Timeout of the context in minutes")
	flagSet.Int64Var(&opts.newSize, "new-size", migrate_volume.USE_EQUAL_SIZE, "Use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC")
	flagSet.StringSliceVar(&opts.nodeSelector, "node-selector", []string{}, "comma separated list of node labels used for nodeSelector of the migration job")
	flagSet.StringVarP(&opts.migrationMode, "migration-mode", "m", "rsync", "Migration mode to use. Options: rsync, rclone")
	flagSet.StringVarP(&opts.migrationFlags, "migration-flags", "f", "", "Additional flags to pass to the migration tool")
}

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewMigrateVolumeCmd(runOptions *migrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = NewMigrateVolumeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the filesystem on a persistent volume to another storage class",
		Long:  `Migrate the filesystem on a persistent volume to another storage class. This will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var cancel context.CancelFunc
			if runOptions.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			}
			defer cancel()
			return migrate_volume.StartMigrateVolumeJob(ctx, migrate_volume.MigrationOptions{
				StorageClassName: runOptions.storageClassName,
				PVCName:          runOptions.pvcName,
				TargetNamespace:  runOptions.targetNamespace,
				NewSize:          runOptions.newSize,
				NodeSelector:     runOptions.nodeSelector,
				MigrationMode:    migrate_volume.MigrationMode(runOptions.migrationMode),
				MigrationFlags:   runOptions.migrationFlags,
			})
		},
	}
	AddMigrateVolumeOptions(cmd.Flags(), runOptions)
	return cmd
}
