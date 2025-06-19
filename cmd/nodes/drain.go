package nodes

import (
	"context"
	_ "embed"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/avisi-cloud/acloud-toolkit/pkg/ame/nodedrain"
)

type maintenanceDrainOptions struct {
	namespace             string
	ignoreStatefulSetPods bool
	uncordonAfterDrain    bool
	gracePeriodSeconds    int
	timeout               time.Duration
}

func newmaintenanceDrainOptions() *maintenanceDrainOptions {
	return &maintenanceDrainOptions{}
}

func AddMaintenanceCreateFlags(flagSet *flag.FlagSet, opts *maintenanceDrainOptions) {
	flagSet.StringVarP(&opts.namespace, "namespace", "n", "", "drain pods from a specific namespace only. Default is the configured namespace in your kubecontext.")
	flagSet.BoolVar(&opts.ignoreStatefulSetPods, "ignore-statefulset-pods", false, "do not drain statefulset pods")
	flagSet.BoolVar(&opts.uncordonAfterDrain, "uncordon", false, "uncordon nodes after running the drain command")

	flagSet.IntVar(&opts.gracePeriodSeconds, "grace-period", 60, "Period of time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used")
	flagSet.DurationVar(&opts.timeout, "timeout", 0*time.Second, "The length of time to wait before giving up, zero means infinite")
}

//go:embed examples/drain.txt
var drainExamples string

// NewDrainCmd returns the Cobra Bootstrap sub command
func NewDrainCmd(runOptions *maintenanceDrainOptions) *cobra.Command {
	if runOptions == nil {
		runOptions = newmaintenanceDrainOptions()
	}

	var cmd = &cobra.Command{
		Use:     "drain <node>",
		Args:    cobra.MinimumNArgs(0),
		Short:   `drain a kubernetes node with additional options not supported by kubectl`,
		Long:    `The acloud-toolkit nodes drain command is a CLI tool that allows you to gracefully remove a Kubernetes node from service, ensuring that all workloads running on the node are rescheduled to other nodes in the cluster before the node is taken offline for maintenance or other purposes. This command provides additional options that are not supported by the standard kubectl drain command.`,
		Example: drainExamples,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			return nodedrain.DrainNode(ctx, args, nodedrain.DrainOptions{
				Namespace:          runOptions.namespace,
				StatelessOnly:      runOptions.ignoreStatefulSetPods,
				UncordonAfterDrain: runOptions.uncordonAfterDrain,
				GracePeriodSeconds: runOptions.gracePeriodSeconds,
				Timeout:            runOptions.timeout,
			})
		},
	}

	AddMaintenanceCreateFlags(cmd.Flags(), runOptions)

	return cmd
}
