package nodedrain

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/drain"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
)

type DrainOptions struct {
	Namespace          string
	StatelessOnly      bool
	UncordonAfterDrain bool
	GracePeriodSeconds int
	Timeout            time.Duration
}

func DrainNode(ctx context.Context, nodeNames []string, opts DrainOptions) error {
	k8sClient := k8s.GetClientOrDie()

	for _, nodeName := range nodeNames {
		if err := cordonNode(ctx, k8sClient, nodeName, opts); err != nil {
			return err
		}
	}

	for _, nodeName := range nodeNames {
		if err := drainNode(ctx, k8sClient, nodeName, opts); err != nil {
			return err
		}
	}

	if opts.UncordonAfterDrain {
		for _, nodeName := range nodeNames {
			if err := uncordonNode(ctx, k8sClient, nodeName, opts); err != nil {
				return err
			}
		}
	}
	return nil
}

func UncordonNode(ctx context.Context, nodeNames []string, opts DrainOptions) error {
	k8sClient := k8s.GetClientOrDie()
	for _, nodeName := range nodeNames {
		if err := uncordonNode(ctx, k8sClient, nodeName, opts); err != nil {
			return err
		}
	}
	return nil
}

func cordonNode(ctx context.Context, k8sClient *kubernetes.Clientset, nodeName string, opts DrainOptions) error {
	fmt.Printf("[acloud-toolkit] cordoning node %q\n", nodeName)
	drainHelper := drain.Helper{
		Ctx:    ctx,
		Client: k8sClient,
	}

	node, err := k8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return drain.RunCordonOrUncordon(&drainHelper, node, true)
}

func uncordonNode(ctx context.Context, k8sClient *kubernetes.Clientset, nodeName string, opts DrainOptions) error {
	fmt.Printf("[acloud-toolkit] uncordoning node %q\n", nodeName)
	drainHelper := drain.Helper{
		Ctx:    ctx,
		Client: k8sClient,
	}

	node, err := k8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return drain.RunCordonOrUncordon(&drainHelper, node, false)
}

func drainNode(ctx context.Context, k8sClient *kubernetes.Clientset, nodeName string, opts DrainOptions) error {
	drainHelper := drain.Helper{
		Ctx:                 ctx,
		Client:              k8sClient,
		Force:               false,
		GracePeriodSeconds:  opts.GracePeriodSeconds,
		IgnoreAllDaemonSets: true,
		Timeout:             opts.Timeout,
		DeleteEmptyDirData:  true,
		Out:                 os.Stdout,
		ErrOut:              os.Stderr,
		// DisableEviction:     drainOpts.DisableEviction,
		OnPodDeletedOrEvicted: func(pod *corev1.Pod, usingEviction bool) {
			if usingEviction {
				fmt.Printf("[acloud-toolkit] evicted pod \"%s/%s\"\n", pod.Namespace, pod.Name)
			} else {
				fmt.Printf("[acloud-toolkit] deleted pod \"%s/%s\"\n", pod.Namespace, pod.Name)
			}
		},
	}

	podsToDelete, errs := drainHelper.GetPodsForDeletion(nodeName)
	if len(errs) != 0 {
		return fmt.Errorf("failed (%d errors) to get pods for deletion for node %q: %w", len(errs), nodeName, errs[0])
	}
	filteredPodsToDelete := []corev1.Pod{}

	for _, pod := range podsToDelete.Pods() {
		if opts.StatelessOnly {
			isStatefulsetPod := false
			for _, ownerReference := range pod.OwnerReferences {
				isStatefulsetPod = isStatefulsetPod || ownerReference.Kind == "StatefulSet"
			}
			if isStatefulsetPod {
				continue
			}
		}
		if opts.Namespace != "" && pod.GetNamespace() == opts.Namespace {
			filteredPodsToDelete = append(filteredPodsToDelete, pod)
		} else if opts.Namespace == "" {
			filteredPodsToDelete = append(filteredPodsToDelete, pod)
		}
	}

	err := drainHelper.DeleteOrEvictPods(filteredPodsToDelete)
	if err != nil {
		return err
	}

	fmt.Printf("[acloud-toolkit] finished draining node %s\n", nodeName)
	return nil
}
