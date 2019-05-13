package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeployments(t *testing.T) {
	// create a dummy deployment
	dummyDeploys := []runtime.Object{
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "myapp"}},
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "prod", Name: "wordpress"}},
	}
	k := KubernetesAPI{
		Client: fake.NewSimpleClientset(dummyDeploys...),
	}
	// Since the Field Selector functionality is done server side, the fake client
	// does not support it. So we can not test with blacklisting here.
	deploymentList, err := k.Deployments([]string{})
	if err != nil {
		t.Error(err)
	}
	deployments := deploymentList.Items

	assert.Len(t, deployments, 2)
}

func TestStatefulSets(t *testing.T) {
	// create a dummy statefulsets
	dummyStatefulSets := []runtime.Object{
		&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "myapp"}},
		&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Namespace: "prod", Name: "wordpress"}},
	}
	k := KubernetesAPI{
		Client: fake.NewSimpleClientset(dummyStatefulSets...),
	}
	// Since the Field Selector functionality is done server side, the fake client
	// does not support it. So we can not test with blacklisting here.
	statefulSetList, err := k.StatefulSets([]string{})
	if err != nil {
		t.Error(err)
	}
	statefulsets := statefulSetList.Items

	assert.Len(t, statefulsets, 2)
}

func TestNodes(t *testing.T) {
	// create a dummy nodes
	dummyNodes := []runtime.Object{
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nodepool-1",
				Labels: map[string]string{
					"name": "nodepool-1",
					"failure-domain.beta.kubernetes.io/region": "westeurope",
					"failure-domain.beta.kubernetes.io/zone":   "1",
					"beta.kubernetes.io/instance-type":         "Standard_D2s_v3",
				},
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nodepool-2",
				Labels: map[string]string{
					"name": "nodepool-2",
					"failure-domain.beta.kubernetes.io/region": "westeurope",
					"failure-domain.beta.kubernetes.io/zone":   "2",
					"beta.kubernetes.io/instance-type":         "Standard_D8s_v3",
				},
			},
		},
	}
	k := KubernetesAPI{
		Client: fake.NewSimpleClientset(dummyNodes...),
	}

	nodes, err := k.Nodes()
	if err != nil {
		t.Error(err)
	}

	assert.Len(t, nodes.Items, 2)
}

func TestBlacklistFieldSelector(t *testing.T) {
	blacklist := []string{"kube-system", "private"}

	fieldSelector := BlacklistFieldSelector(blacklist)

	assert.Equal(t, "metadata.namespace!=kube-system,metadata.namespace!=private", fieldSelector)
}

func TestPrefix(t *testing.T) {
	list := []string{"foo", "bar"}
	prefix := "new-"

	r := Prefix(list, prefix)

	assert.Equal(t, []string{"new-foo", "new-bar"}, r)
}
