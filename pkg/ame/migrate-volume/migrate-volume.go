package migrate_volume

import (
	"context"
	"fmt"
    kubeerrors "k8s.io/apimachinery/pkg/api/errors"
    "time"

	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/helpers"
	"gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func MigrateVolumeJob(ctx context.Context, storageClassName string, pvcName string, namespace string) error {
	k8sClient := k8s.GetClientOrDie()

	migrateVolumeJob := k8sClient.BatchV1().Jobs(namespace)

	jobName := "migrate-volume-" + pvcName
	tmpPVCName := "tmp-" + pvcName

	if err := validateStorageClassExists(ctx, k8sClient, storageClassName); err != nil {
		return err
	}

	pvc, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, pvcName, namespace)
	if err != nil {
		return err
	}

	err = k8s.CreatePersistentVolumeClaim(ctx, k8sClient, tmpPVCName, namespace, storageClassName, *pvc.Spec.Resources.Requests.Storage())
	if err != nil {
		return err
	}
	fmt.Printf("Temporary pvc %q created\n", tmpPVCName)

	ttlSecondsAfterFinished := int32(1000)

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{

				Spec: v1.PodSpec{
					SecurityContext: &v1.PodSecurityContext{
						RunAsNonRoot: helpers.False(),
					},
					Containers: []v1.Container{
						{
							Name:            "volume-migrator",
							Image:           "registry.avisi.cloud/library/rsync:v1",
							ImagePullPolicy: v1.PullAlways,
							Command:         []string{"/bin/sh"},
							Args:            []string{"-c", "rsync -a --stats --progress /mnt/old/ /mnt/new"},
							VolumeMounts: []v1.VolumeMount{
								*k8s.NewVolumeMount("old", "/mnt/old/", true),
								*k8s.NewVolumeMount("new", "/mnt/new/", false),
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
						*k8s.NewPersistentVolumeClaimVolume("old", pvcName, false),
						*k8s.NewPersistentVolumeClaimVolume("new", tmpPVCName, false),
					},
				},
			},
		},
	}

	_, err = migrateVolumeJob.Create(ctx, jobSpec, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job %q: %w", jobName, err)
	}

	err = waitForJobToComplete(ctx, k8sClient, namespace, jobName)
	if err != nil {
		return err
	}

	background := metav1.DeletePropagationBackground
	err = migrateVolumeJob.Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &background})
	if err != nil {
		return fmt.Errorf("failed to delete job %q: %w", jobName, err)
	}

	fmt.Printf("Deleting job: %q\n", jobName)
	time.Sleep(5 * time.Second)
	fmt.Printf("Job %q deleted\n", jobName)

	tmpPVC, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, tmpPVCName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get persistent volume claim%q: %w", tmpPVCName, err)
	}

	err = k8s.SetPVReclaimPolicyToRetain(ctx, k8sClient, pvc)
	if err != nil {
		return err
	}

	err = k8s.SetPVReclaimPolicyToRetain(ctx, k8sClient, tmpPVC)
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

	err = k8s.RemoveClaimRefOfPV(ctx, k8sClient, tmpPVC)
	if err != nil {
		return err
	}

	claimRef := v1.ObjectReference{Name: pvcName, Namespace: namespace}
	err = k8s.SetClaimRefOfPV(ctx, k8sClient, tmpPVC.Spec.VolumeName, claimRef)
	if err != nil {
		return err
	}

	err = k8s.CreatePersistentVolumeClaim(ctx, k8sClient, pvcName, namespace, storageClassName, *pvc.Spec.Resources.Requests.Storage())
	if err != nil {
		return err
	}
	fmt.Printf("Created final pvc %q\n", pvcName)

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

	if *finalPVC.Spec.StorageClassName != storageClassName {
		return fmt.Errorf("new persistent volume claim %q has the storageclass %q and not the given storageclass %q", pvcName, *finalPVC.Spec.StorageClassName, storageClassName)
	}

	fmt.Printf("Data in %q succesfully migrated to %q bound to PVC %q with storageclass %q\n", pvc.Spec.VolumeName, finalPVC.Spec.VolumeName, finalPVC.Name, *finalPVC.Spec.StorageClassName)

	return nil
}

func waitForPVCToBeDeleted(ctx context.Context, k8sClient kubernetes.Interface, namespace, pvc string) error {
	for {
		_, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvc, metav1.GetOptions{})
		if err != nil && kubeerrors.IsNotFound(err)  {
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

func waitForJobToComplete(ctx context.Context, k8sClient kubernetes.Interface, namespace, jobName string) error {
	for {
		job, err := k8sClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if job.Status.Active > 0 {
			fmt.Printf("%s job stil running\n", job.Name)
		}
		if job.Status.Failed > 0 {
			return fmt.Errorf("%s job failed", job.Name)
		}
		if job.Status.Succeeded > 0 {
			fmt.Printf("%s job succeeded\n", job.Name)
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
