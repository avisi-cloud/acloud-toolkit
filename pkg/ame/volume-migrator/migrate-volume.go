package volume_migrator

import (
    "context"
    helpers "gitlab.avisi.cloud/ame/csi-snapshot-utils/pkg/k8s"
    batchv1 "k8s.io/api/batch/v1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "log"
)

func MigrateVolumeJob(clientset *kubernetes.Clientset, jobName *string, pvcNameOld *string, pvcNameNew *string, namespace *string) {
    jobs := clientset.BatchV1().Jobs(*namespace)
    ttlSecondsAfterFinished := int32(1000)

    jobSpec := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      *jobName,
            Namespace: *namespace,
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
                                helpers.NewVolumeMount("old", "/mnt/old/", false),
                                helpers.NewVolumeMount("new", "/mnt/new/", false),
                            },
                        },
                    },
                    RestartPolicy: v1.RestartPolicyNever,
                    Volumes: []v1.Volume{
                        helpers.NewPersistentVolumeClaimVolume("old", *pvcNameOld, false),
                        helpers.NewPersistentVolumeClaimVolume("new", *pvcNameNew, false),
                    },
                },
            },
        },
    }
    _, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
    if err != nil {
        log.Fatalln("Failed to create volume migrator job.")
    }

    //print job details
    log.Println("Created volume migrator job successfully")
}

func SetPVReclaimPolicyToRetain(clientset *kubernetes.Clientset, pv *string){
    persistentVolume, err := clientset.CoreV1().PersistentVolumes().Get(context.TODO(), *pv, metav1.GetOptions{})
    if err != nil {
        log.Fatalln("Failed to get PV with name: " + *pv)
    }

    persistentVolume.Spec.PersistentVolumeReclaimPolicy = "Retain"

    _, err = clientset.CoreV1().PersistentVolumes().Update(context.TODO(), persistentVolume, metav1.UpdateOptions{})
    if err != nil {
        log.Fatalln("Failed to set PVReclaimPolicy")
    }
}



