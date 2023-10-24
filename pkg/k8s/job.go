package k8s

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func WaitForJobToComplete(ctx context.Context, k8sClient kubernetes.Interface, namespace, jobName string) error {
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

func DeleteJobAndWaitForDeletion(ctx context.Context, k8sClient kubernetes.Interface, namespace, jobName string) error {
	background := metav1.DeletePropagationBackground
	err := k8sClient.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &background,
	})
	if err != nil {
		return err
	}
	for {
		_, err := k8sClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
		if err != nil {
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
