package maintenance

import (
	"context"
	"fmt"
	"time"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/helpers"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *NodeUpgraderClient) runMaintenanceJobOnNodeWithScript(ctx context.Context, maintenanceTaskName, nodeName string, scriptContent string, preJobExecHook, postJobExecHook func(node v1.Node)) error {
	scriptName := fmt.Sprintf("acloud-toolkit-%s-%s", maintenanceTaskName, nodeName)
	scriptPath := "task.sh"

	namespace := "kube-system"

	err := k8s.CreateOrUpdateSecret(ctx, c.clusterK8sClient, k8s.CreateSecret(k8s.CreateSecretOptions{
		Name:      scriptName,
		Namespace: namespace,
		Data: map[string][]byte{
			scriptPath: []byte(scriptContent),
		},
	}))
	if err != nil {
		return err
	}
	// make sure we remove the secret once we are done with it
	defer c.clusterK8sClient.CoreV1().Secrets(namespace).Delete(context.Background(), scriptName, metav1.DeleteOptions{})

	return c.runMaintenanceJobOnNode(ctx, maintenanceTaskName, nodeName, scriptName, scriptPath, preJobExecHook, postJobExecHook)
}

func (c *NodeUpgraderClient) runMaintenanceJobOnNode(ctx context.Context, maintenanceTaskName, nodeName string, scriptName, scriptPath string, preJobExecHook, postJobExecHook func(node v1.Node)) error {
	namespace := "kube-system"

	upgradeNodeJob := c.clusterK8sClient.BatchV1().Jobs(namespace)
	ttlSecondsAfterFinished := int32(1000)

	jobName := fmt.Sprintf("%s-%s", maintenanceTaskName, nodeName)

	if err := c.cleanupJob(ctx, jobName); err != nil && !kubeerrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete existing job %q: %w", jobName, err)
	}

	node, err := c.clusterK8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to find node: %q", err)
	}
	if node == nil {
		return fmt.Errorf("no node returned for %s", nodeName)
	}

	if preJobExecHook != nil {
		preJobExecHook(*node)
	}

	_, err = c.clusterK8sClient.CoreV1().Secrets("kube-system").Get(ctx, scriptName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to find node maintenance script")
	}

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            helpers.Int32(4),
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					// Select node based on the hostname label.
					// Note that this is not always equal to the node name, which is why we copy it across from the node object
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": node.Labels["kubernetes.io/hostname"],
					},
					SecurityContext: &v1.PodSecurityContext{
						RunAsNonRoot: helpers.False(),
					},
					HostNetwork:                   true,
					HostPID:                       true,
					TerminationGracePeriodSeconds: helpers.Int64(30),
					Tolerations: []v1.Toleration{
						{
							Key:      "node.kubernetes.io/memory-pressure",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
						{
							Key:      "node.kubernetes.io/disk-pressure",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
						{
							Key:      "node.kubernetes.io/pid-pressure",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
						{
							Key:      "node.kubernetes.io/unschedulable",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
						{
							Key:      "node.kubernetes.io/network-unavailable",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
					},
					Containers: []v1.Container{
						{
							Name:            "node-maintenance",
							Image:           "registry.avisi.cloud/cache/library/centos:7",
							ImagePullPolicy: v1.PullAlways,
							Command:         []string{"chroot", "/host"},
							Args:            []string{"/bin/sh", fmt.Sprintf("/run/maintenance-script/%s", scriptPath)},
							Env:             []v1.EnvVar{},
							VolumeMounts: []v1.VolumeMount{
								k8s.NewVolumeMount("host", "/host", false),
								k8s.NewVolumeMount("maintenance-scripts", "/host/run/maintenance-script", true),
							},
							SecurityContext: &v1.SecurityContext{
								RunAsUser:              helpers.Int64(0),
								RunAsGroup:             helpers.Int64(0),
								ReadOnlyRootFilesystem: helpers.False(),
								Privileged:             helpers.True(),
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
					Volumes: []v1.Volume{
						k8s.NewHostPathVolume("host", "/", v1.HostPathDirectory),
						k8s.NewVolumeFromSecret("maintenance-scripts", scriptName),
					},
				},
			},
		},
	}

	_, err = upgradeNodeJob.Create(ctx, jobSpec, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job %q: %w", jobName, err)
	}
	fmt.Printf("started maintenance job for %q...\n", nodeName)

	// TODO: timeout, 10 minutes enough?
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Should we wait for job to start first?
	err = waitForJobToComplete(timeoutCtx, c.clusterK8sClient, namespace, jobName)
	if err != nil {
		return err
	}
	fmt.Printf("maintenance job for %q has completed\n", nodeName)

	// Clean-up once the job has been completed
	if err := c.cleanupJob(ctx, jobName); err != nil && !kubeerrors.IsNotFound(err) {
		return fmt.Errorf("failed to clean-up upgrade job %q: %w", jobName, err)
	}

	node, err = c.clusterK8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to find node: %q", err)
	}
	if node == nil {
		return fmt.Errorf("no node returned for %s", nodeName)
	}
	if postJobExecHook != nil {
		postJobExecHook(*node)
	}
	return nil
}
