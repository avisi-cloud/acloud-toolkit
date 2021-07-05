package k8s

import v1 "k8s.io/api/core/v1"

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
