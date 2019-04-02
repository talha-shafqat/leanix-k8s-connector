package main

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubernetesAPI is an optionated facade for the kubernetes api
type KubernetesAPI struct {
	Client kubernetes.Interface
}

// NewKubernetesAPI creates a new kuberntes api client
func NewKubernetesAPI(config *rest.Config) (*KubernetesAPI, error) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KubernetesAPI{
		Client: clientset,
	}, nil
}

// Deployments returns a list of deployments filted by the given blacklisted namespaces
func (k *KubernetesAPI) Deployments(blacklistedNamespaces []string) (*appsv1.DeploymentList, error) {
	fieldSelector := BlacklistFieldSelector(blacklistedNamespaces)
	deployments, err := k.Client.AppsV1().Deployments("").List(metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

// Nodes gets the list of worker nodes (kubelets)
func (k *KubernetesAPI) Nodes() (*corev1.NodeList, error) {
	nodes, err := k.Client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// BlacklistFieldSelector builds a Field Selector string to filter the reponse to not
// include resources, that live in the blacklisted namespaces.
func BlacklistFieldSelector(blacklist []string) string {
	namespaceSelectors := Prefix(blacklist, "metadata.namespace!=")
	return strings.Join(namespaceSelectors, ",")
}

// Prefix return a new list where all items are prefixed with the string given as prefix
func Prefix(l []string, p string) []string {
	r := make([]string, 0)
	for _, e := range l {
		r = append(r, (p + e))
	}
	return r
}
