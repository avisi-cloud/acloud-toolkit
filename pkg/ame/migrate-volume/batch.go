package migrate_volume

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"

	v1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BatchMigrateOptions struct {
	SourceStorageClassName string
	TargetStorageClassName string
	TargetNamespace        string
	Timeout                int32
	DryRun                 bool
	MigrationMode          MigrationMode
	MigrationFlags         string
	NodeSelector           []string
}

func BatchMigrateVolumes(ctx context.Context, opts BatchMigrateOptions) error {
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
	namespace := opts.TargetNamespace
	if namespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		namespace = contextNamespace
	}

	pvcs, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volumes in namespace %q: %s", namespace, err)
	}
	if err := k8s.ValidateStorageClassExists(ctx, k8sClient, opts.SourceStorageClassName); err != nil {
		return err
	}
	if err := k8s.ValidateStorageClassExists(ctx, k8sClient, opts.TargetStorageClassName); err != nil {
		return err
	}

	pvcsToMigrate := []string{}
	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase != v1.ClaimBound {
			continue
		}
		if *pvc.Spec.StorageClassName == opts.SourceStorageClassName {
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

	if opts.DryRun {
		fmt.Printf("volumes to migrate: %q (dry-run)\n", pvcsToMigrate)
		return nil
	}

	fmt.Printf("volumes to migrate: %q\n", pvcsToMigrate)
	for _, pvcName := range pvcsToMigrate {
		fmt.Printf("-----\n")
		fmt.Printf("starting volume migration job for PVC \"%s/%s\" ...\n", namespace, pvcName)
		err = StartMigrateVolumeJob(ctx, MigrationOptions{
			StorageClassName: opts.TargetStorageClassName,
			PVCName:          pvcName,
			TargetNamespace:  namespace,
			NewSize:          USE_EQUAL_SIZE,
			MigrationMode:    opts.MigrationMode,
			MigrationFlags:   opts.MigrationFlags,
			NodeSelector:     opts.NodeSelector,
		})

		if err != nil {
			return fmt.Errorf("failed to migrate volume job: %s", err)
		}
		fmt.Printf("finished volume migration job for PVC \"%s/%s\" succesfully\n", namespace, pvcName)
	}

	return nil
}
