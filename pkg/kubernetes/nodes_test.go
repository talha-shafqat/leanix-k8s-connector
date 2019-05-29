package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

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
	k := API{
		Client: fake.NewSimpleClientset(dummyNodes...),
	}

	nodes, err := k.Nodes()
	if err != nil {
		t.Error(err)
	}

	assert.Len(t, nodes.Items, 2)
}
