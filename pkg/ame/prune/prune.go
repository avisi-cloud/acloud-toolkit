package prune

import (
	"context"
	"fmt"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"

	v1 "k8s.io/api/core/v1"
	kresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type Opts struct {
	DryRun        bool
	AllNamespaces bool
	PvcNamespace  string
}

func Volumes(ctx context.Context, opts Opts) error {
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
	namespace := opts.PvcNamespace
	if namespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		namespace = contextNamespace
	}

	pvs, err := k8sClient.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	orphanedPVs := []v1.PersistentVolume{}
	// collect PVs
	for _, pv := range pvs.Items {
		if pv.Status.Phase != v1.VolumeReleased {
			continue
		}
		if pv.Spec.PersistentVolumeReclaimPolicy != v1.PersistentVolumeReclaimDelete {
			orphanedPVs = append(orphanedPVs, pv)
		}
	}

	pvsToPrune := []v1.PersistentVolume{}
	for _, pv := range orphanedPVs {
		if pv.Spec.CSI == nil {
			continue
		}

		if !opts.AllNamespaces {
			if pv.Spec.ClaimRef != nil && pv.Spec.ClaimRef.Namespace != namespace { // filter by PVC namespace
				continue
			}
		}

		if !containsVolumeIdMultipleTimes(pv.Spec.CSI.VolumeHandle, pvs.Items) {
			pvsToPrune = append(pvsToPrune, pv)
		} else if pv.Spec.ClaimRef != nil {
			fmt.Printf("[warning] pruning persistent volume %q (\"%s/%s\")  has volumeHandle that is found multiple times\n", pv.Name, pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
		} else {
			fmt.Printf("[warning] pruning persistent volume %q has volumeHandle that is found multiple times\n", pv.Name)
		}

	}

	if opts.DryRun {
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

func containsVolumeIdMultipleTimes(volumeId string, pvs []v1.PersistentVolume) bool {
	numberOfMatches := 0
	for _, pv := range pvs {
		if pv.Spec.CSI != nil && pv.Spec.CSI.VolumeHandle == volumeId {
			numberOfMatches++

			if numberOfMatches > 1 {
				return true
			}
		}
	}
	return numberOfMatches > 1
}
