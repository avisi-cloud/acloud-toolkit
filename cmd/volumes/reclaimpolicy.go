package volumes

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/reclaimpolicy"
)

type reclaimPolicyOptions struct {
	pvName    string
	pvcName   string
	namespace string
	policy    string
}

func newReclaimPolicyOptions() *reclaimPolicyOptions {
	return &reclaimPolicyOptions{}
}

func AddReclaimPolicyFlags(flagSet *flag.FlagSet, opts *reclaimPolicyOptions) {
	flagSet.StringVar(&opts.pvName, "pv", "", "name of the persistent volume")
	flagSet.StringVar(&opts.pvcName, "pvc", "", "name of the persistent volume claim")
	flagSet.StringVarP(&opts.namespace, "namespace", "n", "", "namespace of the persistent volume claim (optional when using --pvc, defaults to current kubeconfig context)")
	flagSet.StringVarP(&opts.policy, "policy", "p", "", "reclaim policy to set (Retain, Delete, Recycle)")
}

// NewReclaimPolicyCmd returns the Cobra command for changing PV reclaim policy
func NewReclaimPolicyCmd(runOptions *reclaimPolicyOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newReclaimPolicyOptions()
	}

	cmd := &cobra.Command{
		Use:   "reclaim-policy",
		Short: "Change the reclaim policy of a persistent volume",
		Long: `Change the reclaim policy of a persistent volume. You can specify either a persistent volume name directly or a persistent volume claim name.

The reclaim policy determines what happens to the underlying storage when the persistent volume is released:
- Retain: The volume will be retained and must be manually reclaimed
- Delete: The volume will be automatically deleted when released
- Recycle: The volume will be scrubbed and made available for new claims (deprecated)`,
		Example: `
# Set reclaim policy to Retain for a specific PV
acloud-toolkit volumes reclaim-policy --pv my-pv --policy Retain

# Set reclaim policy to Delete for a PV via PVC using current namespace from kubeconfig
acloud-toolkit volumes reclaim-policy --pvc my-pvc --policy Delete

# Set reclaim policy to Retain for a PV via PVC in a specific namespace
acloud-toolkit volumes reclaim-policy --pvc data-pvc --namespace production --policy Retain
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runOptions.pvName == "" && runOptions.pvcName == "" {
				return fmt.Errorf("either --pv or --pvc must be specified")
			}
			if runOptions.pvName != "" && runOptions.pvcName != "" {
				return fmt.Errorf("cannot specify both --pv and --pvc")
			}

			return reclaimpolicy.SetReclaimPolicy(cmd.Context(), reclaimpolicy.ReclaimPolicyOptions{
				PVName:    runOptions.pvName,
				PVCName:   runOptions.pvcName,
				Namespace: runOptions.namespace,
				Policy:    runOptions.policy,
			})
		},
	}

	AddReclaimPolicyFlags(cmd.Flags(), runOptions)

	cmd.MarkFlagRequired("policy")

	return cmd
}
