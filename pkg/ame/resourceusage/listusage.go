package resourceusage

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/k8s"
	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/table"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func List(sourceNamespace string, allNamespaces bool) error {
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

	namespace, _, err := kubeconfig.Namespace()
	if err != nil {
		return err
	}
	if namespace == "" {
		return nil
	}
	if sourceNamespace != "" {
		namespace = sourceNamespace
	}

	// Collect namespaces, either based a single (default) namespace, or collect all namespaces in the cluster
	// and use those to query the snapshots in cluster.
	namespaces := []string{}
	if allNamespaces {
		k8sclient := k8s.GetClientOrDie()
		namespaceList, err := k8sclient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, ns := range namespaceList.Items {
			namespaces = append(namespaces, ns.GetName())
		}
	} else {
		namespaces = append(namespaces, namespace)
	}

	// collect all deployments for each namespace
	deployments := []appsv1.Deployment{}
	statefulsets := []appsv1.StatefulSet{}
	for _, namespace := range namespaces {
		deploymentList, err := k8sClient.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		deployments = append(deployments, deploymentList.Items...)

		stsList, err := k8sClient.AppsV1().StatefulSets(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		statefulsets = append(statefulsets, stsList.Items...)
	}

	// Format output
	body := make([][]string, 0, len(deployments)+len(statefulsets))
	for _, deployment := range deployments {
		totalMemoryPerPod := resource.NewQuantity(0, resource.BinarySI)
		for _, container := range deployment.Spec.Template.Spec.Containers {
			memory := container.Resources.Limits.Memory()
			if memory != nil {
				totalMemoryPerPod.Add(*memory)
			}
		}

		replicas := deployment.Spec.Replicas

		body = append(body, []string{
			deployment.GetNamespace(),
			"Deployment",
			deployment.Name,
			fmt.Sprintf("%d", len(deployment.Spec.Template.Spec.Containers)),
			fmt.Sprintf("%d", *replicas),
			totalMemoryPerPod.String(),
		})
	}

	for _, sts := range statefulsets {
		totalMemoryPerPod := resource.NewQuantity(0, resource.BinarySI)
		for _, container := range sts.Spec.Template.Spec.Containers {
			memory := container.Resources.Limits.Memory()
			if memory != nil {
				totalMemoryPerPod.Add(*memory)
			}
		}

		replicas := sts.Spec.Replicas

		body = append(body, []string{
			sts.GetNamespace(),
			"StatefulSet",
			sts.Name,
			fmt.Sprintf("%d", len(sts.Spec.Template.Spec.Containers)),
			fmt.Sprintf("%d", *replicas),
			totalMemoryPerPod.String(),
		})
	}

	table.Print([]string{
		"Namespace",
		"Type",
		"Name",
		"Containers",
		"Replicas",
		"Memory",
	}, body)

	return nil
}
