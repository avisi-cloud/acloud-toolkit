package storage

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
}

func NewMigrateVolumeOptions() *migrateVolumeOptions {
	return &migrateVolumeOptions{}
}

func AddMigrateVolumeOptions(flagSet *flag.FlagSet, opts *migrateVolumeOptions) {
	flagSet.StringVarP(&opts.storageClassName, "storageClass", "s", "", "name of the new storageclass")
	flagSet.StringVarP(&opts.pvcName, "pvc", "p", "", "name of the persitentvolumeclaim")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "default", "Namespace where de migrate job will be executed")
	flagSet.Int32VarP(&opts.timeout, "timeout", "t", 60, "Timeout of the context in minutes")
	flagSet.Int64Var(&opts.newSize, "new-size", migrate_volume.USE_EQUAL_SIZE, "Use a different size for the new PVC. Value is in MB. Default 0 means use same size as current PVC")
}

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewMigrateVolumeCmd(ctx context.Context, runOptions *migrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = NewMigrateVolumeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate a volume to another storage class",
		Long:  `Migrate a volume to another storage class. This will create a new PVC using the target storage class, and copy all file contents over to the new volume. The existing persistent volume will remain available in the cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			defer cancel()
			return migrate_volume.MigrateVolumeJob(ctxWithTimeout, runOptions.storageClassName, runOptions.pvcName, runOptions.targetNamespace, runOptions.newSize)
		},
	}

	AddMigrateVolumeOptions(cmd.Flags(), runOptions)

	return cmd
}
