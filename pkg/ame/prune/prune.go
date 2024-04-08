package prune

import (
	"context"
	"fmt"
	"time"

	volumesnapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	v1 "k8s.io/api/core/v1"
	kresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
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

		if !containsVolumeIdMultipleTimes(pv.Spec.CSI.VolumeHandle, pvs.Items) {
			pvsToPrune = append(pvsToPrune, pv)
		} else if pv.Spec.ClaimRef != nil {
			fmt.Printf("[warning] pruning persistent volume %q (\"%s/%s\")  has volumeHandle that is found multiple times\n", pv.Name, pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
		} else {
			fmt.Printf("[warning] pruning persistent volume %q has volumeHandle that is found multiple times\n", pv.Name)
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

func PruneSnapshots(ctx context.Context, namespace string, allNamespaces, pruneImported, dryRun bool, minAge time.Duration) error {
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
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	if namespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		namespace = contextNamespace
	}

	namespaces := []string{}
	if allNamespaces {
		namespaceList, err := k8sClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, ns := range namespaceList.Items {
			namespaces = append(namespaces, ns.GetName())
		}
	} else {
		namespaces = append(namespaces, namespace)
	}

	// Collect snapshots for each namespace
	var listSnapshots []volumesnapshotv1.VolumeSnapshot
	for _, namespace := range namespaces {
		snapshotUnstructured, err := dynamicClient.Resource(volumesnapshotv1.SchemeGroupVersion.WithResource("volumesnapshots")).Namespace(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		if snapshotUnstructured == nil {
			continue
		}
		for _, snapshotItem := range snapshotUnstructured.Items {
			snapshot := volumesnapshotv1.VolumeSnapshot{}

			err := runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotItem.Object, &snapshot)
			if err != nil {
				return err
			}

			listSnapshots = append(listSnapshots, snapshot)
		}
	}

	// Figure out which snapshots to prune
	var filteredSnapshots []volumesnapshotv1.VolumeSnapshot
	for _, snapshot := range listSnapshots {

		// Skip snapshots that have the app.kubernetes.io/managed-by annotation set
		if snapshot.Annotations["app.kubernetes.io/managed-by"] != "" {
			continue
		}

		// Skip snapshots that are not ready to use
		if snapshot.Status.ReadyToUse == nil || !*snapshot.Status.ReadyToUse {
			continue
		}

		volumeSnapshotContent, err := getBoundVolumeSnapshotContentForSnapshot(ctx, snapshot, dynamicClient)
		if err != nil {
			return err
		}
		if volumeSnapshotContent == nil {
			fmt.Printf("skipping snapshot %q as it has no bound volume snapshot content\n", snapshot.Name)
			continue
		}

		// Skip snapshots that have their deletion policy set to anything other than delete
		if volumeSnapshotContent.Spec.DeletionPolicy != volumesnapshotv1.VolumeSnapshotContentDelete {
			continue
		}

		// Skip snapshots that are not older than minAge
		if minAge > 0 {
			age := time.Since(snapshot.CreationTimestamp.Time)
			if age < minAge {
				continue
			}
		}

		filteredSnapshots = append(filteredSnapshots, snapshot)
	}

	fmt.Println("Snapshots to prune:")
	for _, snapshot := range filteredSnapshots {
		if dryRun {
			fmt.Printf("pruning snapshot %q ... (dry run)\n", snapshot.Name)
		} else {
			fmt.Printf("pruning snapshot %q ...\n", snapshot.Name)
		}
	}

	return nil
}

func getBoundVolumeSnapshotContentForSnapshot(ctx context.Context, snapshot volumesnapshotv1.VolumeSnapshot, client dynamic.Interface) (*volumesnapshotv1.VolumeSnapshotContent, error) {
	snapshotContentName := snapshot.Status.BoundVolumeSnapshotContentName
	if snapshotContentName == nil {
		return nil, nil
	}

	snapshotContentUnstructured, err := client.Resource(volumesnapshotv1.SchemeGroupVersion.WithResource("volumesnapshotcontents")).Get(ctx, *snapshotContentName, metav1.GetOptions{})
	if err != nil {
		return &volumesnapshotv1.VolumeSnapshotContent{}, err
	}
	snapshotContent := volumesnapshotv1.VolumeSnapshotContent{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(snapshotContentUnstructured.Object, &snapshotContent)
	if err != nil {
		return &volumesnapshotv1.VolumeSnapshotContent{}, err
	}

	return &snapshotContent, nil
}
