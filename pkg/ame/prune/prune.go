package prune

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"

	v1 "k8s.io/api/core/v1"
	kresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func PruneVolumes(ctx context.Context, dryRun bool) error {
	kubeconfig, err := k8s.GetClientConfig()
	if err != nil {
		return err
	}
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return err
	}
	k8sClient, err := k8s.GetClientWithConfig(config)
	if err != nil {
		return err
	}

	pvs, err := k8sClient.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	pvsToPrune := []v1.PersistentVolume{}
	// collect PVs
	for _, pv := range pvs.Items {
		if pv.Status.Phase != v1.VolumeReleased {
			continue
		}
		if pv.Spec.PersistentVolumeReclaimPolicy != v1.PersistentVolumeReclaimDelete {
			pvsToPrune = append(pvsToPrune, pv)
		}
	}

	if dryRun {
		totalSize := int64(0)
		for _, pv := range pvsToPrune {
			storageCapacity := pv.Spec.Capacity.Storage()
			if storageCapacity != nil {
				asInt, _ := storageCapacity.AsInt64()
				totalSize += asInt
			}
			if pv.Spec.ClaimRef != nil {
				fmt.Printf("pruning persistent volume %q (\"%s/%s\") ... (dry run)\n", pv.Name, pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
			} else {
				fmt.Printf("pruning persistent volume %q ... (dry run)\n", pv.Name)
			}
		}
		q := kresource.NewQuantity(totalSize, kresource.BinarySI)
		fmt.Printf("total storage pruned: %s (dry run)\n", q.String())
		return nil
	}

	totalSize := int64(0)
	for _, pv := range pvsToPrune {
		storageCapacity := pv.Spec.Capacity.Storage()
		if storageCapacity != nil {
			asInt, _ := storageCapacity.AsInt64()
			totalSize += asInt
		}
		if pv.Spec.ClaimRef != nil {
			fmt.Printf("pruning persistent volume %q (\"%s/%s\") ...\n", pv.Name, pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
		} else {
			fmt.Printf("pruning persistent volume %q ...\n", pv.Name)
		}

		pv.Spec.PersistentVolumeReclaimPolicy = v1.PersistentVolumeReclaimDelete
		_, err = k8sClient.CoreV1().PersistentVolumes().Update(ctx, &pv, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update reclaim policy persistent volume %q: %w", pv.Name, err)
		}
	}
	q := kresource.NewQuantity(totalSize, kresource.BinarySI)
	fmt.Printf("total storage pruned: %s\n", q.String())
	return nil
}
