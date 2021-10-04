package storage

import (
	"context"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/resize"
)

type volumeResizeOptions struct {
	namespace string
	name      string
	newSize   string
}

func newvolumeResizeOptions() *volumeResizeOptions {
	return &volumeResizeOptions{}
}

func AddvolumeResizeFlags(flagSet *flag.FlagSet, opts *volumeResizeOptions) {
	flagSet.StringVar(&opts.newSize, "size", "", "New size. Example: 10G")
	flagSet.StringVarP(&opts.namespace, "namespace", "n", "default", "Namespace of the PVC. Snapshot will be created within this namespace as well")
	flagSet.StringVarP(&opts.name, "pvc", "p", "", "Name of the persistent volume to snapshot")
}

// NewvolumeResizeCmd returns the Cobra Bootstrap sub command
func NewvolumeResizeCmd(runOptions *volumeResizeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newvolumeResizeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "resize <persistent-volume-claim>",
		Short: "resize adjusts the volume size of the pvc",
		Long:  `resize adjusts the volume size of the pvc`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if err := resize.ResizeVolume(context.Background(), runOptions.namespace, arg, runOptions.newSize); err != nil {
					return err
				}
			}
			return nil
		},
	}

	AddvolumeResizeFlags(cmd.Flags(), runOptions)

	return cmd
}
