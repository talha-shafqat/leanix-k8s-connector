package kubernetes

import (
	"github.com/leanix/leanix-k8s-connector/pkg/set"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSets returns a list of statefulsets filted by the given blacklisted namespaces
func (k *API) StatefulSets(blacklistedNamespaces []string) (*appsv1.StatefulSetList, error) {
	fieldSelector := BlacklistFieldSelector(blacklistedNamespaces)
	statefulsets, err := k.Client.AppsV1().StatefulSets("").List(metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return statefulsets, nil
}

// StatefulSetsOnNodes returns a list of statefulSets filted by the given blacklisted namespaces
// and the nodes the statefulSet pods are running on
func (k *API) StatefulSetsOnNodes(blacklistedNamespaces []string) (*appsv1.StatefulSetList, *[]corev1.Node, error) {
	statefulSets, err := k.StatefulSets(blacklistedNamespaces)
	if err != nil {
		return nil, nil, err
	}
	nodeNames := set.NewStringSet()
	for _, s := range statefulSets.Items {
		pods, err := k.StatefulSetPods(&s)
		if err != nil {
			return nil, nil, err
		}
		for _, p := range pods.Items {
			nodeNames.Add(p.Spec.NodeName)
		}
	}
	nodes, err := k.NodesByName(nodeNames)
	return statefulSets, nodes, nil
}

// StatefulSetPods retuns a list of pods matching the selectors of the given statefulSet
func (k *API) StatefulSetPods(statefulSet *appsv1.StatefulSet) (*corev1.PodList, error) {
	return k.Pods(statefulSet.Spec.Selector.MatchLabels)
}
