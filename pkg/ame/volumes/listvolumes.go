package volumes

import (
	"context"
	"fmt"

	"github.com/avisi-cloud/acloud-toolkit/pkg/k8s"

	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type PersistentVolumeInfo struct {
	PeristentVolume v1.PersistentVolume
	Claim           *v1.PersistentVolumeClaim
	Attachment      *storagev1.VolumeAttachment
}

func ListVolumes(ctx context.Context, listUnattachedOnly bool, filterByStorageClassName string) ([]PersistentVolumeInfo, error) {
	kubeconfig, err := k8s.GetClientConfig()
	if err != nil {
		return nil, err
	}
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := k8s.GetClientWithConfig(config)
	if err != nil {
		return nil, err
	}

	volumeInfo := []PersistentVolumeInfo{}

	pvs, err := k8sClient.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	attachments, err := k8sClient.StorageV1().VolumeAttachments().List(ctx, metav1.ListOptions{})
	if err != nil && !kubeerrors.IsNotFound(err) {
		return nil, fmt.Errorf("error while checking for volume attachments: %s", err)
	}

	for _, pv := range pvs.Items {
		if filterByStorageClassName != "" && pv.Spec.StorageClassName != filterByStorageClassName {
			continue
		}
		var volumeAttachment *storagev1.VolumeAttachment
		for _, attachment := range attachments.Items {
			if attachment.Spec.Source.PersistentVolumeName != nil && *attachment.Spec.Source.PersistentVolumeName == pv.Name {
				volumeAttachment = &attachment
				break
			}
		}

		if listUnattachedOnly && volumeAttachment != nil && volumeAttachment.Status.Attached {
			continue
		}

		var pvc *v1.PersistentVolumeClaim
		if pv.Spec.ClaimRef != nil {
			pvc, err = k8sClient.CoreV1().PersistentVolumeClaims(pv.Spec.ClaimRef.Namespace).Get(ctx, pv.Spec.ClaimRef.Name, metav1.GetOptions{})
			if err != nil && !kubeerrors.IsNotFound(err) {
				return nil, err
			}
		}

		volumeInfo = append(volumeInfo, PersistentVolumeInfo{
			PeristentVolume: pv,
			Attachment:      volumeAttachment,
			Claim:           pvc,
		})
	}
	return volumeInfo, nil
}
