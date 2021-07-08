package migrate_volume

import (
	"context"
	"fmt"
	"time"

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

	pvc, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, pvcName, namespace)
	if err != nil {
		return err
	}

	err = k8s.CreatePersistentVolumeClaim(ctx, k8sClient, tmpPVCName, namespace, storageClassName, *pvc.Spec.Resources.Requests.Storage())
	if err != nil {
		return err
	}
	fmt.Printf("Temporary pvc %s created\n", tmpPVCName)

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
					Containers: []v1.Container{
						{
							Name:    "volume-migrator",
							Image:   "centos:7",
							Command: []string{"/bin/sh"},
							Args:    []string{"-c", "yum -y install rsync","rsync -r /mnt/old/ /mnt/new"},
							VolumeMounts: []v1.VolumeMount{
								*k8s.NewVolumeMount("old", "/mnt/old/", true),
								*k8s.NewVolumeMount("new", "/mnt/new/", false),
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

	fmt.Printf("Delete job: %s\n", jobName)
	time.Sleep(5 * time.Second)
	fmt.Printf("Job: %s deleted\n", jobName)

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
	fmt.Printf("Deleting pvc: %s\n", tmpPVCName)

	err = k8sClient.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim%q: %w", pvcName, err)
	}
	fmt.Printf("Deleting pvc: %s\n", pvcName)

	err = k8s.RemoveClaimRefOfPV(ctx, k8sClient, tmpPVC)
	if err != nil {
		return err
	}

	err = k8s.CreatePersistentVolumeClaim(ctx, k8sClient, pvcName, namespace, storageClassName, *pvc.Spec.Resources.Requests.Storage())
	if err != nil {
		return err
	}
	fmt.Printf("Creating pvc: %s\n", pvcName)

	claimRef := v1.ObjectReference{Name: pvcName, Namespace: namespace}
	err = k8s.SetClaimRefOfPV(ctx, k8sClient, tmpPVC.Spec.VolumeName, claimRef)
	if err != nil {
		return err
	}

	return nil
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
			return fmt.Errorf("%s job failed\n", job.Name)
		}
		if job.Status.Succeeded > 0 {
			fmt.Printf("%s job succeeded\n", job.Name)
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			continue
		}
	}
}
