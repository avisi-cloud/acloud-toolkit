package k8s

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeApply struct {
	client    kubernetes.Interface
	namespace string
}

func NewKubeApply(namespace string) (*KubeApply, error) {
	k8sclient, err := GetClient()
	if err != nil {
		return nil, err
	}

	return &KubeApply{
		client:    k8sclient,
		namespace: namespace,
	}, nil
}

func NewKubeApplyOrDie(namespace string) *KubeApply {
	client, err := NewKubeApply(namespace)
	if err != nil {
		panic(err)
	}
	return client
}

func NewKubeApplyWithClient(namespace string, client kubernetes.Interface) *KubeApply {
	return &KubeApply{
		client:    client,
		namespace: namespace,
	}
}

func (a *KubeApply) ApplyStatefulSet(ctx context.Context, sts *appsv1.StatefulSet) error {
	_, err := a.client.AppsV1().StatefulSets(a.namespace).Create(ctx, sts, metav1.CreateOptions{})
	if err != nil {
		if kubeerrors.IsAlreadyExists(err) {
			_, err := a.client.AppsV1().StatefulSets(a.namespace).Update(ctx, sts, metav1.UpdateOptions{})
			return err
		}
		return err
	}
	return nil
}

func (a *KubeApply) ApplySecret(ctx context.Context, secret *v1.Secret) error {
	err := CreateOrUpdateSecret(ctx, a.client, secret)
	if err != nil {
		return err
	}
	return nil
}

func (a *KubeApply) ApplyConfigMap(ctx context.Context, configmap *v1.ConfigMap) error {
	err := CreateOrUpdateConfigMap(ctx, a.client, configmap)
	if err != nil {
		return err
	}
	return nil
}

func (a *KubeApply) ApplyDeployment(ctx context.Context, deploy *appsv1.Deployment) error {
	_, err := a.client.AppsV1().Deployments(a.namespace).Create(ctx, deploy, metav1.CreateOptions{})
	if err != nil {
		if kubeerrors.IsAlreadyExists(err) {
			_, err := a.client.AppsV1().Deployments(a.namespace).Update(ctx, deploy, metav1.UpdateOptions{})
			return err
		}
		return err
	}
	return nil
}

func (a *KubeApply) DeleteDeployment(ctx context.Context, name string) error {
	err := a.client.AppsV1().Deployments(a.namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if !kubeerrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (a *KubeApply) ScaleDeployment(ctx context.Context, deploymentName string, replicas int32) error {
	scale, err := a.client.AppsV1().Deployments(a.namespace).GetScale(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = replicas
	_, err = a.client.AppsV1().Deployments(a.namespace).UpdateScale(ctx, deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a *KubeApply) ScaleStatefulSet(ctx context.Context, statefulSetName string, replicas int32) error {
	scale, err := a.client.AppsV1().StatefulSets(a.namespace).GetScale(ctx, statefulSetName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = replicas
	_, err = a.client.AppsV1().StatefulSets(a.namespace).UpdateScale(ctx, statefulSetName, scale, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a *KubeApply) ApplyService(ctx context.Context, svc *v1.Service) error {
	existingSvc, err := a.client.CoreV1().Services(a.namespace).Get(ctx, svc.GetName(), metav1.GetOptions{})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			_, err := a.client.CoreV1().Services(a.namespace).Create(ctx, svc, metav1.CreateOptions{})
			if err != nil && !kubeerrors.IsAlreadyExists(err) {
				return errors.Wrapf(err, fmt.Sprintf("failed install service %s", svc.GetName()))
			}
			return nil
		}
		return errors.Wrapf(err, fmt.Sprintf("failed install service %s", svc.GetName()))
	}

	if existingSvc == nil {
		_, err := a.client.CoreV1().Services(a.namespace).Create(ctx, svc, metav1.CreateOptions{})
		if err != nil && !kubeerrors.IsAlreadyExists(err) {
			return errors.Wrapf(err, fmt.Sprintf("failed install service %s", svc.GetName()))
		}
		return nil
	}

	existingSvc.Spec.Ports = svc.Spec.Ports
	existingSvc.Spec.SessionAffinity = svc.Spec.SessionAffinity
	existingSvc.Spec.Selector = svc.Spec.Selector
	existingSvc.Spec.PublishNotReadyAddresses = svc.Spec.PublishNotReadyAddresses
	existingSvc.Spec.ExternalTrafficPolicy = svc.Spec.ExternalTrafficPolicy
	existingSvc.Labels = svc.Labels
	existingSvc.Annotations = svc.Annotations

	_, err = a.client.CoreV1().Services(a.namespace).Update(ctx, existingSvc, metav1.UpdateOptions{})
	return err
}

func (a *KubeApply) ApplyIngress(ctx context.Context, ingress *networkingv1.Ingress) error {
	const errorMsg = "failed to install ingress"

	existingIngress, err := a.client.NetworkingV1().Ingresses(a.namespace).Get(ctx, ingress.GetName(), metav1.GetOptions{})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			_, err := a.client.NetworkingV1().Ingresses(a.namespace).Create(ctx, ingress, metav1.CreateOptions{})
			if err != nil && !kubeerrors.IsAlreadyExists(err) {
				return errors.Wrapf(err, errorMsg)
			}
			return nil
		}
		return errors.Wrapf(err, errorMsg)
	}

	if existingIngress == nil {
		return fmt.Errorf(errorMsg)
	}

	// Update the attributes
	existingIngress.Spec.DefaultBackend = ingress.Spec.DefaultBackend
	existingIngress.Spec.TLS = ingress.Spec.TLS
	existingIngress.Spec.Rules = ingress.Spec.Rules
	if ingress.Spec.IngressClassName != nil {
		existingIngress.Spec.IngressClassName = ingress.Spec.IngressClassName
	}
	existingIngress.Labels = ingress.Labels
	existingIngress.Annotations = ingress.Annotations

	// Patch the installed ingress with out new values
	_, err = a.client.NetworkingV1().Ingresses(a.namespace).Update(ctx, existingIngress, metav1.UpdateOptions{})
	return err
}

func (a *KubeApply) ApplyNetworkPolicy(ctx context.Context, networkPolicy networkingv1.NetworkPolicy) error {
	const errorMsg = "failed to install network policy"

	existingPSP, err := a.client.NetworkingV1().NetworkPolicies(a.namespace).Get(ctx, networkPolicy.GetName(), metav1.GetOptions{})
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			_, err := a.client.NetworkingV1().NetworkPolicies(a.namespace).Create(ctx, &networkPolicy, metav1.CreateOptions{})
			if err != nil && !kubeerrors.IsAlreadyExists(err) {
				return errors.Wrapf(err, errorMsg)
			}
			return nil
		}
		return errors.Wrapf(err, errorMsg)
	}

	if existingPSP == nil {
		return fmt.Errorf(errorMsg)
	}

	// Update the attributes
	existingPSP.Spec = networkPolicy.Spec
	existingPSP.Labels = networkPolicy.Labels
	existingPSP.Annotations = networkPolicy.Annotations

	// Patch the installed network policy with the new values
	_, err = a.client.NetworkingV1().NetworkPolicies(a.namespace).Update(ctx, &networkPolicy, metav1.UpdateOptions{})
	return err
}
