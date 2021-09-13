package migrate_volume

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/k8s"

	v1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func validateStorageClassExists(ctx context.Context, client *kubernetes.Clientset, storageClassName string) error {
	_, err := client.StorageV1().StorageClasses().Get(ctx, storageClassName, metav1.GetOptions{})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			return fmt.Errorf("storage class %q does not exist", storageClassName)
		}
		return fmt.Errorf("error while checking storage class: %s", err)
	}
	return nil
}

func BatchMigrateVolumes(ctx context.Context, sourceStorageClass, targetStorageClass, namespace string, dryRun bool) error {
	k8sClient := k8s.GetClientOrDie()
	pvcs, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volumes in namespace %q: %s", namespace, err)
	}
	if err := validateStorageClassExists(ctx, k8sClient, sourceStorageClass); err != nil {
		return err
	}
	if err := validateStorageClassExists(ctx, k8sClient, targetStorageClass); err != nil {
		return err
	}

	pvcsToMigrate := []string{}
	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase != v1.ClaimBound {
			continue
		}
		if *pvc.Spec.StorageClassName == sourceStorageClass {
			// check if the volume is attached

			// if no volumeName, this PVC has no bound PV and should not have to be migrated
			if pvc.Spec.VolumeName == "" {
				continue
			}

			attachments, err := k8sClient.StorageV1().VolumeAttachments().List(ctx, metav1.ListOptions{})

			isAttached := false
			for _, attachment := range attachments.Items {
				if attachment.Spec.Source.PersistentVolumeName != nil && *attachment.Spec.Source.PersistentVolumeName == pvc.Spec.VolumeName {
					if err != nil && !kubeerrors.IsNotFound(err) {
						return fmt.Errorf("error while checking for volume attachments: %s", err)
					}
					if attachment.Status.Attached {

						isAttached = true
						fmt.Printf("volume %s is still attached. Skipping ...\n", pvc.GetName())
					}
					break
				}
			}
			if !isAttached {
				pvcsToMigrate = append(pvcsToMigrate, pvc.GetName())
			}
		}
	}

	if dryRun {
		fmt.Printf("volumes to migrate: %q (dry-run)\n", pvcsToMigrate)
		return nil
	}
	fmt.Printf("volumes to migrate: %q\n", pvcsToMigrate)
	for _, pvcName := range pvcsToMigrate {
		fmt.Printf("-----\n")
		fmt.Printf("starting volume migration job for PVC \"%s/%s\" ...\n", namespace, pvcName)
		err = MigrateVolumeJob(ctx, targetStorageClass, pvcName, namespace, USE_EQUAL_SIZE)
		if err != nil {
			return fmt.Errorf("failed to migrate volume job: %s", err)
		}
		fmt.Printf("finished volume migration job for PVC \"%s/%s\" succesfully\n", namespace, pvcName)
	}

	return nil
}
