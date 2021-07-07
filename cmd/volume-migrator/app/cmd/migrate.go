package cmd

import (
    "context"
    "github.com/spf13/cobra"
    flag "github.com/spf13/pflag"
    "gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/ame/migrate-volume"
)

type migrateVolumeOptions struct {
    storageClassName      string
    PVCname      string
    targetNamespace string
}


func NewMigrateVolumeOptions() *migrateVolumeOptions {
    return &migrateVolumeOptions{}
}

func AddMigrateVolumeOptions(flagSet *flag.FlagSet, opts *migrateVolumeOptions) {
    flagSet.StringVarP(&opts.storageClassName, "storageClass", "s", "", "name of the new storageclass")
    flagSet.StringVarP(&opts.PVCname, "pvc", "p", "", "name of the persitentvolumeclaim")
    flagSet.StringVarP(&opts.targetNamespace, "target-namespace","n", "default", "Namespace where de migrate job will be executed")
}

// NewMigrateVolumeCmd returns the Cobra Bootstrap sub command
func NewMigrateVolumeCmd(ctx context.Context,runOptions *migrateVolumeOptions) *cobra.Command {
    if runOptions == nil {
        runOptions = NewMigrateVolumeOptions()
    }

    var cmd = &cobra.Command{
        Use:   "migrate",
        Short: "Migrate a volume",
        Long:  `Migrate a volume from one PVC to other PVC`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return migrate_volume.MigrateVolumeJob(ctx, runOptions.storageClassName, runOptions.PVCname, runOptions.targetNamespace)
        },
    }

    AddMigrateVolumeOptions(cmd.Flags(), runOptions)

    return cmd
}

