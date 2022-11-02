package maintenance

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	upgradeSecretSecretName = "avisi-cloud-node-upgrade-script"
)

type NodeUpgraderClient struct {
	// clusterK8sClient holds a reference to the clusterSettings's Kubernetes API server.
	// Should not be nil.
	clusterK8sClient clientset.Interface
}

type Opts struct {
}

// NewNodeMaintenanceClient creates a new client for running maintenance tasks on kubernetes nodes
func NewNodeMaintenanceClient(clusterK8sClient clientset.Interface, opts Opts) *NodeUpgraderClient {
	return &NodeUpgraderClient{
		clusterK8sClient: clusterK8sClient,
	}
}

func (c *NodeUpgraderClient) cleanupJob(ctx context.Context, jobName string) error {
	namespace := "kube-system"
	upgradeNodeJob := c.clusterK8sClient.BatchV1().Jobs(namespace)

	background := metav1.DeletePropagationBackground
	err := upgradeNodeJob.Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &background})
	if err != nil && !kubeerrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete existing upgrade job %q: %w", jobName, err)
	}
	return nil
}

func waitForJobToComplete(ctx context.Context, k8sClient clientset.Interface, namespace, jobName string) error {
	for {
		job, err := k8sClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if job.Status.Succeeded > 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
			continue
		}
	}
}
