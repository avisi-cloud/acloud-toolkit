package reclaimpolicy

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
)

type ReclaimPolicyOptions struct {
	PVName    string
	PVCName   string
	Namespace string
	Policy    string
}

func SetReclaimPolicy(ctx context.Context, opts ReclaimPolicyOptions) error {
	k8sClient, err := k8s.GetClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	var policy v1.PersistentVolumeReclaimPolicy
	switch opts.Policy {
	case "Retain":
		policy = v1.PersistentVolumeReclaimRetain
	case "Delete":
		policy = v1.PersistentVolumeReclaimDelete
	case "Recycle":
		policy = v1.PersistentVolumeReclaimRecycle
	default:
		return fmt.Errorf("invalid reclaim policy %q, must be one of: Retain, Delete, Recycle", opts.Policy)
	}

	if opts.PVName != "" {
		return setPVReclaimPolicyByName(ctx, k8sClient, opts.PVName, policy)
	}

	if opts.PVCName != "" {
		namespace := opts.Namespace
		if namespace == "" {
			kubeconfig, err := k8s.GetClientConfig()
			if err != nil {
				return fmt.Errorf("failed to get kubernetes config: %w", err)
			}
			contextNamespace, _, err := kubeconfig.Namespace()
			if err != nil {
				return fmt.Errorf("failed to get namespace from kubeconfig: %w", err)
			}
			namespace = contextNamespace
		}
		return setPVReclaimPolicyByPVC(ctx, k8sClient, opts.PVCName, namespace, policy)
	}

	return fmt.Errorf("either --pv or --pvc must be specified")
}

func setPVReclaimPolicyByName(ctx context.Context, k8sClient kubernetes.Interface, pvName string, policy v1.PersistentVolumeReclaimPolicy) error {
	pv, err := k8sClient.CoreV1().PersistentVolumes().Get(ctx, pvName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volume %q: %w", pvName, err)
	}

	if pv.Spec.PersistentVolumeReclaimPolicy == policy {
		fmt.Printf("PV %s already has %s as the reclaim policy\n", pvName, policy)
		return nil
	}

	fmt.Printf("Updating PV %s reclaim policy from %s to %s...\n", pvName, pv.Spec.PersistentVolumeReclaimPolicy, policy)
	pv.Spec.PersistentVolumeReclaimPolicy = policy
	_, err = k8sClient.CoreV1().PersistentVolumes().Update(ctx, pv, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update reclaim policy for persistent volume %q: %w", pvName, err)
	}

	fmt.Printf("Successfully updated PV %s reclaim policy to %s\n", pvName, policy)
	return nil
}

func setPVReclaimPolicyByPVC(ctx context.Context, k8sClient kubernetes.Interface, pvcName, namespace string, policy v1.PersistentVolumeReclaimPolicy) error {
	pvc, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volume claim %q in namespace %q: %w", pvcName, namespace, err)
	}

	if pvc.Spec.VolumeName == "" {
		return fmt.Errorf("persistent volume claim %q is not bound to a volume", pvcName)
	}

	return setPVReclaimPolicyByName(ctx, k8sClient, pvc.Spec.VolumeName, policy)
}
