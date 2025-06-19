package syncvolume

import (
	"context"
	"fmt"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/avisi-cloud/acloud-toolkit/pkg/helpers"
	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"
)

const (
	USE_EQUAL_SIZE = 0
)

type SyncVolumeJobOptions struct {
	Namespace               string
	SourcePVCName           string
	TargetPVCName           string
	NewStorageClassName     string
	CreateNewPVC            bool
	RetainJob               bool
	ExtraRsyncArgs          []string
	TtlSecondsAfterFinished int32
	NewSize                 int64
}

func SyncVolumeJob(ctx context.Context, opts SyncVolumeJobOptions) error {
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
	if opts.Namespace == "" {
		contextNamespace, _, err := kubeconfig.Namespace()
		if err != nil {
			return err
		}
		opts.Namespace = contextNamespace
	}

	syncVolumeJob := k8sClient.BatchV1().Jobs(opts.Namespace)
	jobName := helpers.FormatKubernetesName(fmt.Sprintf("sync-"+opts.SourcePVCName+"-to-"+opts.TargetPVCName), helpers.MaxKubernetesLabelValueLength, 5)

	sourcePVC, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, opts.SourcePVCName, opts.Namespace)
	if err != nil {
		return fmt.Errorf("failed to get source-pvc %q: %w", opts.SourcePVCName, err)
	}

	if opts.CreateNewPVC {
		if err := k8s.ValidateStorageClassExists(ctx, k8sClient, opts.NewStorageClassName); err != nil {
			return fmt.Errorf("storage class %q does not exist: %w", opts.NewStorageClassName, err)
		}

		sourcePVCStorageSize := *sourcePVC.Spec.Resources.Requests.Storage()
		if opts.NewSize > USE_EQUAL_SIZE {
			sourcePVCStorageSize = resource.MustParse(fmt.Sprintf("%dM", opts.NewSize))
		}

		err = k8s.CreatePersistentVolumeClaim(ctx, k8sClient, opts.TargetPVCName, opts.Namespace, opts.NewStorageClassName, sourcePVCStorageSize)
		if err != nil {
			// if the pvc already exists while create-pvc option is true, an existing pvc is not used
			if kubeerrors.IsAlreadyExists(err) {
				return fmt.Errorf("target PVC %q already exists, to use an existing PVC remove the --create-pvc option", opts.TargetPVCName)
			}
			return err
		}
		fmt.Printf("pvc %q created\n", opts.TargetPVCName)
	}

	if _, err := k8s.GetPersistentVolumeClaimAndCheckForVolumes(ctx, k8sClient, opts.TargetPVCName, opts.Namespace); err != nil {
		return fmt.Errorf("failed to get target-pvc %q: %w", opts.TargetPVCName, err)
	}

	rsyncCmd := fmt.Sprintf("rsync -a --stats --progress /mnt/source/ /mnt/target/ %s", strings.Join(opts.ExtraRsyncArgs, " "))

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: opts.Namespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &opts.TtlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{

				Spec: v1.PodSpec{
					SecurityContext: &v1.PodSecurityContext{
						RunAsNonRoot: helpers.False(),
					},
					Containers: []v1.Container{
						{
							Name:            "volume-sync",
							Image:           "registry.avisi.cloud/library/rsync:v1",
							ImagePullPolicy: v1.PullAlways,
							Command:         []string{"/bin/sh"},
							Args:            []string{"-c", rsyncCmd},
							VolumeMounts: []v1.VolumeMount{
								k8s.NewVolumeMount("source", "/mnt/source/", true),
								k8s.NewVolumeMount("target", "/mnt/target/", false),
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
						k8s.NewPersistentVolumeClaimVolume("source", opts.SourcePVCName, false),
						k8s.NewPersistentVolumeClaimVolume("target", opts.TargetPVCName, false),
					},
				},
			},
		},
	}

	if _, err := syncVolumeJob.Create(ctx, jobSpec, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create job %q: %w", jobName, err)
	}

	if err := k8s.WaitForJobToComplete(ctx, k8sClient, opts.Namespace, jobName); err != nil {
		return err
	}

	if !opts.RetainJob {
		fmt.Printf("deleting job %q\n", jobName)
		if err := k8s.DeleteJobAndWaitForDeletion(ctx, k8sClient, opts.Namespace, jobName); err != nil {
			return fmt.Errorf("failed to delete job %q: %w", jobName, err)
		}
		fmt.Printf("job %q deleted\n", jobName)
	}

	return nil
}
