package migrate_volume

import (
    "context"
    "fmt"
    "gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/k8s"
    batchv1 "k8s.io/api/batch/v1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MigrateVolumeJob(pvcNameOld string, pvcNameNew string, namespace string) error {
    kubeconfig, err := k8s.GetClientCmd()
    if err != nil {
        return err
    }
    config, err := kubeconfig.ClientConfig()
    if err != nil {
        return err
    }
    k8sclient, err := k8s.GetClientWithConfig(config)
    if err != nil {
        return err
    }

    migrateVolumeJob := k8sclient.BatchV1().Jobs(namespace)

    jobName := "migrate-volume-" + pvcNameOld

    pvc, err := k8s.GetPersistentVolumeClaim(k8sclient, pvcNameOld, namespace)
    if err != nil {
        return err
    }

    err = k8s.SetPVReclaimPolicyToRetain(k8sclient, *pvc)
    if err != nil {
        return err
    }

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
                            Args: []string{"-c", "cp -rp /mnt/old/ /mnt/new/"},
                            VolumeMounts: []v1.VolumeMount{
                                k8s.NewVolumeMount("old", "/mnt/old/", false),
                                k8s.NewVolumeMount("new", "/mnt/new/", false),
                            },
                        },
                    },
                    RestartPolicy: v1.RestartPolicyNever,
                    Volumes: []v1.Volume{
                        k8s.NewPersistentVolumeClaimVolume("old", pvcNameOld, false),
                        k8s.NewPersistentVolumeClaimVolume("new", pvcNameNew, false),
                    },
                },
            },
        },
    }

    _, err = migrateVolumeJob.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
    if err != nil {
        return err
    }

    err = k8s.RemoveClaimRefOfPV(k8sclient, *pvc)
    if err != nil {
        return err
    }

    //print job details
    fmt.Printf("Created volume migrator job successfully")
    return nil
}



