package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDeployments(t *testing.T) {
	// create a dummy deployment
	dummyDeploys := []runtime.Object{
		&v1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "myapp"}},
		&v1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: "prod", Name: "wordpress"}},
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
