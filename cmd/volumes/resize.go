package volumes

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/resize"
)

type volumeResizeOptions struct {
	namespace string
	newSize   string
}

func newvolumeResizeOptions() *volumeResizeOptions {
	return &volumeResizeOptions{}
}

func AddvolumeResizeFlags(flagSet *flag.FlagSet, opts *volumeResizeOptions) {
	flagSet.StringVar(&opts.newSize, "size", "", "New size. Example: 10G")
	flagSet.StringVarP(&opts.namespace, "namespace", "n", "", "If present, the namespace scope for this CLI request. Otherwise uses the namespace from the current Kubernetes context")
}

// NewvolumeResizeCmd returns the Cobra Bootstrap sub command
func NewvolumeResizeCmd(runOptions *volumeResizeOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newvolumeResizeOptions()
	}

	var cmd = &cobra.Command{
		Use:   "resize <persistent-volume-claim>",
		Short: "Resize adjusts the volume size of a persistent volume claim",
		Long:  `The 'resize' command adjusts the size of a persistent volume claim (PVC). The command takes a PVC name as input along with an optional namespace parameter and a new size in gigabytes.`,
		Example: `
# Resize a PVC named 'data' in the default namespace to 20 gigabytes
acloud-toolkit storage resize data --size 20G

# Resize a PVC named 'data' in the 'prod' namespace to 50 gigabytes
acloud-toolkit storage resize data --namespace prod --size 50G	  
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if err := resize.ResizeVolume(cmd.Context(), runOptions.namespace, arg, runOptions.newSize); err != nil {
					return err
				}
			}
			return nil
		},
	}

	AddvolumeResizeFlags(cmd.Flags(), runOptions)

	return cmd
}
