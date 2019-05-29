package kubernetes

import (
	"github.com/leanix/leanix-k8s-connector/pkg/set"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Deployments returns a list of deployments filted by the given blacklisted namespaces
func (k *API) Deployments(blacklistedNamespaces []string) (*appsv1.DeploymentList, error) {
	fieldSelector := BlacklistFieldSelector(blacklistedNamespaces)
	deployments, err := k.Client.AppsV1().Deployments("").List(metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

// DeploymentsOnNodes returns a list of deployments filted by the given blacklisted namespaces
// and the nodes the deploment pods are running on
func (k *API) DeploymentsOnNodes(blacklistedNamespaces []string) (*appsv1.DeploymentList, *[]corev1.Node, error) {
	deployments, err := k.Deployments(blacklistedNamespaces)
	if err != nil {
		return nil, nil, err
	}
	nodeNames := set.NewStringSet()
	for _, d := range deployments.Items {
		pods, err := k.DeploymentPods(&d)
		if err != nil {
			return nil, nil, err
		}
		for _, p := range pods.Items {
			nodeNames.Add(p.Spec.NodeName)
		}
	}
	nodes, err := k.NodesByName(nodeNames)
	return deployments, nodes, nil
}

// DeploymentPods retuns a list of pods matching the selectors of the given deployment
func (k *API) DeploymentPods(deployment *appsv1.Deployment) (*corev1.PodList, error) {
	return k.Pods(deployment.Spec.Selector.MatchLabels)
}
