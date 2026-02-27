package volumes

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/reclaimpolicy"
	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
)

type reclaimPolicyOptions struct {
	pvNames   []string
	pvcNames  []string
	namespace string
	policy    string
}

func newReclaimPolicyOptions() *reclaimPolicyOptions {
	return &reclaimPolicyOptions{}
}

func AddReclaimPolicyFlags(flagSet *flag.FlagSet, opts *reclaimPolicyOptions) {
	flagSet.StringArrayVar(&opts.pvNames, "pv", []string{}, "name of the persistent volume (may be specified multiple times)")
	flagSet.StringArrayVar(&opts.pvcNames, "pvc", []string{}, "name of the persistent volume claim (may be specified multiple times)")
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

# Set reclaim policy to Retain for multiple PVs at once
acloud-toolkit volumes reclaim-policy --pv my-pv --pv another-pv --policy Retain

# Set reclaim policy to Delete for a PV via PVC using current namespace from kubeconfig
acloud-toolkit volumes reclaim-policy --pvc my-pvc --policy Delete

# Set reclaim policy to Retain for multiple PVCs in a specific namespace
acloud-toolkit volumes reclaim-policy --pvc data-pvc --pvc logs-pvc --namespace production --policy Retain
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, pvName := range runOptions.pvNames {
				if err := reclaimpolicy.SetReclaimPolicy(cmd.Context(), reclaimpolicy.ReclaimPolicyOptions{
					PVName: pvName,
					Policy: runOptions.policy,
				}); err != nil {
					return err
				}
			}
			for _, pvcName := range runOptions.pvcNames {
				if err := reclaimpolicy.SetReclaimPolicy(cmd.Context(), reclaimpolicy.ReclaimPolicyOptions{
					PVCName:   pvcName,
					Namespace: runOptions.namespace,
					Policy:    runOptions.policy,
				}); err != nil {
					return err
				}
			}
			return nil
		},
	}

	AddReclaimPolicyFlags(cmd.Flags(), runOptions)

	cmd.MarkFlagRequired("policy")
	cmd.MarkFlagsOneRequired("pv", "pvc")

	_ = cmd.RegisterFlagCompletionFunc("pv", completePVNames)
	_ = cmd.RegisterFlagCompletionFunc("pvc", completePVCNames)

	return cmd
}

// completePVNames returns the names of all PersistentVolumes in the cluster for shell completion.
func completePVNames(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	pvList, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	names := make([]string, 0, len(pvList.Items))
	for _, pv := range pvList.Items {
		if strings.HasPrefix(pv.Name, toComplete) {
			names = append(names, pv.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// completePVCNames returns the names of PersistentVolumeClaims for shell completion.
// It reads the --namespace flag from cmd; if unset it falls back to the current kubeconfig context
// namespace. If the namespace cannot be determined, completion fails with an error.
func completePVCNames(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		clientConfig, err := k8s.GetClientConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		if ns, _, err := clientConfig.Namespace(); err == nil && ns != "" {
			namespace = ns
		}
	}

	client, err := k8s.GetClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	pvcList, err := client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	names := make([]string, 0, len(pvcList.Items))
	for _, pvc := range pvcList.Items {
		if strings.HasPrefix(pvc.Name, toComplete) {
			names = append(names, pvc.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
