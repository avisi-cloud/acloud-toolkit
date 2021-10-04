package resize

import (
	"context"
	"fmt"
	"time"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/helpers"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const (
	retryCount      = 3
	backOffDuration = 2 * time.Second
)

func ResizeVolume(ctx context.Context, namespace string, pvcName string, newSize string) error {
	k8sClient := k8s.GetClientOrDie()

	pvc, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, pvcName, namespace)
	if err != nil {
		return err
	}

	storageSize := *pvc.Spec.Resources.Requests.Storage()
	newStorageSize, err := resource.ParseQuantity(newSize)
	if err != nil {
		return fmt.Errorf("cannot parse size into quantity: %v", err)
	}

	fmt.Printf("current size: %d\n", storageSize.MilliValue())
	fmt.Printf("new size: %d\n", newStorageSize.MilliValue())
	if storageSize.MilliValue() > newStorageSize.MilliValue() {
		return fmt.Errorf("volume can only be expanded")
	}

	fmt.Printf("current size in cluster: %s\n", storageSize.String())
	fmt.Printf("new size in cluster: %s\n", newStorageSize.String())
	// if newStorageSize != storageSize {
	// 	return nil
	// }
	pvc.Spec.Resources.Requests = corev1.ResourceList{
		corev1.ResourceStorage: newStorageSize,
	}

	err = helpers.RetryWithCancel(ctx, retryCount, backOffDuration, func() error {
		_, err = k8sClient.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update persistent volume claim %q: %w", pvc.Name, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = helpers.RetryWithCancel(ctx, retryCount, backOffDuration, func() error {
		// if the PVC is not mounted, we should only wait until the PV has been expanded, since the filesystem cannot be expanded until it's mounted.
		return waitForPVCToBeExpanded(ctx, k8sClient, namespace, pvcName, newStorageSize)
	})
	if err != nil {
		return err
	}
	fmt.Printf("succesfully expanded persistent volume %q\n", pvc.Name)
	return nil
}

func waitForPVCToBeExpanded(ctx context.Context, k8sClient kubernetes.Interface, namespace, pvcName string, newStorageSize resource.Quantity) error {
	for {
		pvc, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error waiting for pvc to expand: %s", err)
		}

		pv, err := k8sClient.CoreV1().PersistentVolumes().Get(ctx, pvc.Spec.VolumeName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error waiting for pvc to expand: %s", err)
		}

		if pv.Spec.Capacity.Storage().Equal(newStorageSize) {
			return nil
		}
		// if pvc.Status.Capacity.Storage().Equal(newStorageSize) {
		// 	return nil
		// }
		fmt.Printf("persistent volume claim %q still in the proces of being expanded...\n", pvc.Name)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			continue
		}
	}
}
