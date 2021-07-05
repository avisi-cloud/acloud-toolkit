package k8s

import (
    "context"
    "fmt"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "time"
)

func NewVolumeMount(name, path string, readOnly bool) v1.VolumeMount {
    return v1.VolumeMount{
        Name:      name,
        MountPath: path,
        ReadOnly:  readOnly,
    }
}

func NewPersistentVolumeClaimVolume(name, claimName string, readOnly bool) v1.Volume {
    return v1.Volume{
        Name: name,
        VolumeSource: v1.VolumeSource{
            PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
                ClaimName: claimName,
                ReadOnly:  readOnly,
            },
        },
    }
}

func SetPVReclaimPolicyToRetain(k8sclient *kubernetes.Clientset, pvc v1.PersistentVolumeClaim) error {
    // Get the persistent volume, ensure it's set to Retain.
    pv, err := k8sclient.CoreV1().PersistentVolumes().Get(context.TODO(), pvc.Spec.VolumeName, metav1.GetOptions{})
    if err != nil {
        return err
    }

    if pv.Spec.PersistentVolumeReclaimPolicy != v1.PersistentVolumeReclaimRetain {
        fmt.Printf("PV %s does not have retain as the reclaim policy, updating ...\n", pvc.Spec.VolumeName)
        pv.Spec.PersistentVolumeReclaimPolicy = v1.PersistentVolumeReclaimRetain
        _, err2 := k8sclient.CoreV1().PersistentVolumes().Update(context.TODO(), pv, metav1.UpdateOptions{})
        if err2 != nil {
            return err2
        }

        // give kube some time to catch up
        time.Sleep(1 * time.Second)
    }

    // give kube some time to catch up
    time.Sleep(1 * time.Second)
    return nil
}

func GetPersistentVolumeClaim(k8sclient *kubernetes.Clientset, pvcName string, namespace string) (*v1.PersistentVolumeClaim,error){
    pvc, err := k8sclient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }
    for {
        if pvc.Spec.VolumeName != "" {
            break
        }
        time.Sleep(1 * time.Second)

        pvc, err = k8sclient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
        if err != nil {
            return nil, err
        }
    }
    return pvc, err
}

func RemoveClaimRefOfPV(k8sclient *kubernetes.Clientset, pvc v1.PersistentVolumeClaim) error{
    // Update the PVC to remove the claimRef
    pv, err := k8sclient.CoreV1().PersistentVolumes().Get(context.TODO(), pvc.Spec.VolumeName, metav1.GetOptions{})
    if err != nil {
        return err
    }
    pv.Spec.ClaimRef = nil
    _, err = k8sclient.CoreV1().PersistentVolumes().Update(context.TODO(), pv, metav1.UpdateOptions{})
    if err != nil {
        return err
    }
    fmt.Printf("removed the PV %s claim ref to %s...\n", pvc.Spec.VolumeName, pvc.Name)
    return nil
}
