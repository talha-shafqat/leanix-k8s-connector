package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestDeployments(t *testing.T) {
	// create a dummy deployment
	dummyDeploys := []runtime.Object{
		&v1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "myapp"}},
	}
	k := KubernetesAPI{
		Client: testclient.NewSimpleClientset(dummyDeploys...),
	}
	blacklist := []string{"kube-system"}
	deploymentList, err := k.Deployments(blacklist)
	if err != nil {
		t.Error(err)
	}
	deployments := deploymentList.Items

	assert.Len(t, deployments, 1)
}

func TestDeployments_FiltersBlacklistedNamespaces(t *testing.T) {
	// create a dummy deployment
	dummyDeploys := []runtime.Object{
		&v1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "myapp"}},
		&v1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "kube-system", Name: "system-app"}},
	}
	k := KubernetesAPI{
		Client: testclient.NewSimpleClientset(dummyDeploys...),
	}
	blacklist := []string{"kube-system"}
	deploymentList, err := k.Deployments(blacklist)
	if err != nil {
		t.Error(err)
	}
	deployments := deploymentList.Items

	assert.Len(t, deployments, 1)
	assert.Equal(t, deployments[0].Name, "myapp")
}

func TestPrefix(t *testing.T) {
	list := []string{"foo", "bar"}
	prefix := "new-"

	r := Prefix(list, prefix)

	assert.Equal(t, []string{"new-foo", "new-bar"}, r)
}
