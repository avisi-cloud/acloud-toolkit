package cmd

import (
	"context"
	"time"

	migrate_volume "gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/ame/migrate-volume"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type migrateVolumeOptions struct {
	storageClassName string
	pvcName          string
	targetNamespace  string
	timeout          int32
}

func NewMigrateVolumeOptions() *migrateVolumeOptions {
	return &migrateVolumeOptions{}
}

func AddMigrateVolumeOptions(flagSet *flag.FlagSet, opts *migrateVolumeOptions) {
	flagSet.StringVarP(&opts.storageClassName, "storageClass", "s", "", "name of the new storageclass")
	flagSet.StringVarP(&opts.pvcName, "pvc", "p", "", "name of the persitentvolumeclaim")
	flagSet.StringVarP(&opts.targetNamespace, "target-namespace", "n", "default", "Namespace where de migrate job will be executed")
	flagSet.Int32VarP(&opts.timeout, "timeout", "t", 60, "Timeout of the context in minutes")
}

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewMigrateVolumeCmd(ctx context.Context, runOptions *migrateVolumeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = NewMigrateVolumeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate a volume",
		Long:  `Migrate a volume from one PVC to other PVC`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(runOptions.timeout)*time.Minute)
			defer cancel()
			return migrate_volume.MigrateVolumeJob(ctxWithTimeout, runOptions.storageClassName, runOptions.pvcName, runOptions.targetNamespace)
		},
	}

	AddMigrateVolumeOptions(cmd.Flags(), runOptions)

	return cmd
}
