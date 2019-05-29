package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
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
	k := API{
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
