package volumes

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/volumes"
	"github.com/avisi-cloud/acloud-toolkit/pkg/table"
	"github.com/avisi-cloud/acloud-toolkit/pkg/timeformat"
)

type listOptions struct {
	storageClassName string
	unattachedOnly   bool
	csiOnly          bool
}

func newListOptions() *listOptions {
	return &listOptions{}
}

func AddListFlags(flagSet *flag.FlagSet, opts *listOptions) {
	flagSet.StringVarP(&opts.storageClassName, "storage-class", "s", "", "run for storage class. Will use default storage class if left empty")
	flagSet.BoolVar(&opts.unattachedOnly, "unattached-only", false, "show unattached persistent volumes only")
	flagSet.BoolVar(&opts.csiOnly, "csi-only", false, "show CSI persistent volumes only")
}

//go:embed examples/list.txt
var listExamples string

// NewListCmd returns the Cobra Bootstrap sub command
func NewListCmd(runOptions *listOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newListOptions()
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all persistent volumes in a Kubernetes cluster",
		Long:    `This command lists all CSI persistent volumes within the cluster. This command allows you to list and filter persistent volumes based on various criteria, making it easier to inspect and manage your storage resources.`,
		Aliases: []string{"ls"},
		Example: listExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			volumes, err := volumes.ListVolumes(cmd.Context(), runOptions.unattachedOnly, runOptions.storageClassName)
			if err != nil {
				return err
			}

			// Format output
			body := make([][]string, 0, len(volumes))
			for _, volume := range volumes {
				size := ""
				if volume.PeristentVolume.Spec.Capacity.Storage() != nil {
					size = volume.PeristentVolume.Spec.Capacity.Storage().String()
				}

				claimNamespace := "-"
				claimName := "-"
				claimStatus := "-"

				if volume.Claim != nil {
					claimNamespace = volume.Claim.GetNamespace()
					claimName = volume.Claim.Name
					claimStatus = string(volume.PeristentVolume.Status.Phase)
				}

				attachmentNode := "-"
				attached := "false" // we cannot discover attachment status from non CSI volumes
				attachmentAge := "-"
				if volume.Attachment != nil {
					attachmentNode = volume.Attachment.Spec.NodeName
					attached = fmt.Sprint(volume.Attachment.Status.Attached)
					attachmentAge = timeformat.FormatTime(volume.Attachment.CreationTimestamp.Time, false)
				}

				isCSI := volume.PeristentVolume.Spec.CSI != nil
				if runOptions.csiOnly && !isCSI {
					continue
				}

				if runOptions.unattachedOnly && attached == "true" {
					continue
				}
				if !isCSI && attached == "false" {
					attached = "-"
				}

				volumeHandle := ""
				if isCSI {
					volumeHandle = volume.PeristentVolume.Spec.CSI.VolumeHandle
				}

				body = append(body, []string{
					volume.PeristentVolume.Name,
					volume.PeristentVolume.Spec.StorageClassName,
					volumeHandle,
					claimNamespace,
					claimName,
					claimStatus,
					size,
					timeformat.FormatTime(volume.PeristentVolume.CreationTimestamp.Time, false),
					attachmentNode,
					fmt.Sprint(attached),
					attachmentAge,
					fmt.Sprint(isCSI),
				})
			}

			table.Print([]string{
				"VolumeName",
				"Storage ClassName",
				"Volume Handle",
				"Claim Namespace",
				"Claim Name",
				"Status",
				"Capacity",
				"Volume Age",
				"Node",
				"Attached",
				"Attachment Age",
				"CSI",
			}, body)

			return nil
		},
	}

	AddListFlags(cmd.Flags(), runOptions)

	return cmd
}
