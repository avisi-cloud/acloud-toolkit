package migrate_volume

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/avisi-cloud/acloud-toolkit/pkg/helpers"
	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const (
	USE_EQUAL_SIZE = 0
)

// MigrationMode is the type of tool used for migration of the filesystem
type MigrationOptions struct {
	// StorageClassName is the name of the new storageclass
	StorageClassName string
	// PVCName is the name of the persistent volume claim
	PVCName string
	// TargetNamespace is the namespace where the volume migrate job will be executed
	TargetNamespace string
	// NewSize is the size of the new PVC. Value is in MB. Default 0 means use same size as current PVC
	NewSize int64

	// MigrationMode is the type of tool used for migration of the filesystem
	MigrationMode MigrationMode

	// RCloneImage is the image used for the rclone migration tool
	RCloneImage string
	// RyncImage is the image used for the rsync migration tool
	RyncImage string

	// comma separated list of node labels used for nodeSelector
	NodeSelector []string

	// Additional flags to pass to the migration tool
	MigrationFlags string
}

// StartMigrateVolumeJob starts a job that migrates the filesystem on a persistent volume to another storage class
// This will create a new PVC using the target storage class, and copy all file contents over to the new volume.
// The existing persistent volume will remain available in the cluster.
// The function will return an error if the migration fails.
func StartMigrateVolumeJob(ctx context.Context, opts MigrationOptions) error {

	metav1.FormatLabelSelector(nil)
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

	migrateVolumeJob := k8sClient.BatchV1().Jobs(namespace)

	pvcName := opts.PVCName
	jobName := helpers.FormatKubernetesName(fmt.Sprintf("migrate-volume-%s", pvcName), helpers.MaxKubernetesLabelValueLength, 5)
	tmpPVCName := "tmp-" + opts.PVCName

	if err := k8s.ValidateStorageClassExists(ctx, k8sClient, opts.StorageClassName); err != nil {
		return err
	}

	pvc, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, pvcName, namespace)
	if err != nil {
		return err
	}

	storageSize := *pvc.Spec.Resources.Requests.Storage()
	if opts.NewSize > USE_EQUAL_SIZE {
		storageSize = resource.MustParse(fmt.Sprintf("%dM", opts.NewSize))
	}

	err = k8s.CreatePersistentVolumeClaim(ctx, k8sClient, tmpPVCName, namespace, opts.StorageClassName, storageSize)
	if err != nil {
		if !kubeerrors.IsAlreadyExists(err) {
			return err
		}
		fmt.Printf("Using existing pvc %q\n", tmpPVCName)
	} else {
		fmt.Printf("Temporary pvc %q created\n", tmpPVCName)
	}

	ttlSecondsAfterFinished := int32(1000)

	image := ""
	args := ""
	switch opts.MigrationMode {
	case MigrationModeRsync:
		image = DefaultRSyncContainerImage
		if opts.RyncImage != "" {
			image = opts.RyncImage
		}
		args = fmt.Sprintf("rsync -a --stats --progress %s /mnt/old/ /mnt/new", opts.MigrationFlags)
	case MigrationModeRclone:
		image = DefaultRCloneContainerImage
		if opts.RCloneImage != "" {
			image = opts.RCloneImage
		}
		args = fmt.Sprintf("rclone copy /mnt/old/ /mnt/new --progress %s", opts.MigrationFlags)
	default:
		return fmt.Errorf("unknown mode %q", opts.MigrationMode)
	}

	labelSelector, err := metav1.ParseToLabelSelector(strings.Join(opts.NodeSelector, ","))
	if err != nil {
		return fmt.Errorf("failed to parse nodeSelector: %w", err)
	}
	nodeSelector, err := metav1.LabelSelectorAsMap(labelSelector)
	if err != nil {
		return fmt.Errorf("failed to convert nodeSelector to map: %w", err)
	}

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					NodeSelector: nodeSelector,
					SecurityContext: &v1.PodSecurityContext{
						RunAsNonRoot: helpers.False(),
					},
					Containers: []v1.Container{
						{
							Name:            "volume-migrator",
							Image:           image,
							ImagePullPolicy: v1.PullAlways,
							Command:         []string{"/bin/sh"},
							Args:            []string{"-c", args},
							VolumeMounts: []v1.VolumeMount{
								k8s.NewVolumeMount("old", "/mnt/old/", true),
								k8s.NewVolumeMount("new", "/mnt/new/", false),
							},
							SecurityContext: &v1.SecurityContext{
								RunAsUser:              helpers.Int64(0),
								RunAsGroup:             helpers.Int64(0),
								ReadOnlyRootFilesystem: helpers.False(),
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
					Volumes: []v1.Volume{
						k8s.NewPersistentVolumeClaimVolume("old", pvcName, false),
						k8s.NewPersistentVolumeClaimVolume("new", tmpPVCName, false),
					},
				},
			},
		},
	}

	_, err = migrateVolumeJob.Create(ctx, jobSpec, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job %q: %w", jobName, err)
	}

	err = k8s.WaitForJobToComplete(ctx, k8sClient, namespace, jobName)
	if err != nil {
		return err
	}

	fmt.Printf("deleting job %q\n", jobName)
	if err := k8s.DeleteJobAndWaitForDeletion(ctx, k8sClient, namespace, jobName); err != nil {
		return fmt.Errorf("failed to delete job %q: %w", jobName, err)
	}
	fmt.Printf("job %q deleted\n", jobName)

	tmpPVC, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, tmpPVCName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volume claim%q: %w", tmpPVCName, err)
	}

	err = helpers.RetryWithCancel(ctx, 3, 2*time.Second, func() error {
		err = k8s.SetPVReclaimPolicyToRetain(ctx, k8sClient, pvc)
		if err != nil {
			return err
		}
		return k8s.SetPVReclaimPolicyToRetain(ctx, k8sClient, tmpPVC)
	})
	if err != nil {
		return err
	}

	err = k8sClient.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, tmpPVCName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim%q: %w", tmpPVCName, err)
	}
	fmt.Printf("Deleting temp pvc %q (persistent volume is marked as retain)\n", tmpPVCName)

	err = k8sClient.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim%q: %w", pvcName, err)
	}
	fmt.Printf("Deleting source pvc: %s (persistent volume is marked as retain)\n", pvcName)

	err = waitForPVCToBeDeleted(ctx, k8sClient, namespace, pvcName)
	if err != nil {
		return err
	}

	err = helpers.RetryWithCancel(ctx, 3, 2*time.Second, func() error {
		err = k8s.RemoveClaimRefOfPV(ctx, k8sClient, tmpPVC)
		if err != nil {
			return err
		}

		claimRef := v1.ObjectReference{Name: pvcName, Namespace: namespace}
		return k8s.SetClaimRefOfPV(ctx, k8sClient, tmpPVC.Spec.VolumeName, claimRef)
	})
	if err != nil {
		return err
	}

	err = helpers.RetryWithCancel(ctx, 3, 2*time.Second, func() error {
		return k8s.CreatePersistentVolumeClaim(ctx, k8sClient, pvcName, namespace, opts.StorageClassName, storageSize)
	})
	if err != nil {
		return err
	}
	fmt.Printf("Created final pvc %q\n", pvcName)

	return helpers.RetryWithCancel(ctx, 3, 2*time.Second, func() error {
		finalPVC, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, pvcName, namespace)
		if err != nil {
			return fmt.Errorf("failed to get new persistent volume claim%q: %w", pvcName, err)
		}

		if finalPVC.Status.Phase != v1.ClaimBound {
			return fmt.Errorf("new persistent volume claim is not bound! %q", tmpPVC.Name)
		}

		if finalPVC.Spec.VolumeName != tmpPVC.Spec.VolumeName {
			return fmt.Errorf("new persistent volume claim %q is not bound to the new persistentvolume! %q", finalPVC.Name, tmpPVC.Spec.VolumeName)
		}

		if *finalPVC.Spec.StorageClassName != opts.StorageClassName {
			return fmt.Errorf("new persistent volume claim %q has the storageclass %q and not the given storageclass %q", pvcName, *finalPVC.Spec.StorageClassName, opts.StorageClassName)
		}
		fmt.Printf("Data in %q succesfully migrated to %q bound to PVC %q with storageclass %q\n", pvc.Spec.VolumeName, finalPVC.Spec.VolumeName, finalPVC.Name, *finalPVC.Spec.StorageClassName)
		return nil
	})
}

func waitForPVCToBeDeleted(ctx context.Context, k8sClient kubernetes.Interface, namespace, pvc string) error {
	for {
		_, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvc, metav1.GetOptions{})
		if err != nil && kubeerrors.IsNotFound(err) {
			fmt.Printf("source pvc %s is deleted\n", pvc)
			return nil
		} else if err != nil {
			return fmt.Errorf("error deleting source pvc: %s", err)
		}
		fmt.Printf("source pvc %s still in the proces of being deleted...\n", pvc)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			continue
		}
	}
}
