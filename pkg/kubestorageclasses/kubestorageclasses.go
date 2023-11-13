package kubestorageclasses

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	ErrNoVolumeSnapshotClassFound   = errors.New("no default volume snapshot class found")
	ErrNoVolumeSnapshotContentFound = errors.New("no default volume snapshot content found")
)

const (
	DefaultStorageClassAnnotation = "storageclass.kubernetes.io/is-default-class"
)

func GetDefaultStorageClassName(ctx context.Context, k8sclient *kubernetes.Clientset) (string, error) {
	storageClasses, err := k8sclient.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("error validating storage classes: %w", err)
	}
	if len(storageClasses.Items) == 0 {
		return "", errors.New("no storage classes present in this cluster")
	}

	for _, sc := range storageClasses.Items {
		value, ok := sc.Annotations[DefaultStorageClassAnnotation]
		if ok && value == "true" {
			return sc.Name, nil
		}
	}
	return "", errors.New("no default storage class installed in cluster")
}
