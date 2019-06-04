package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMapNodes(t *testing.T) {
	clusterName := "mycluster"
	// create a dummy nodes
	nodes := &corev1.NodeList{
		Items: []corev1.Node{
			corev1.Node{
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
			corev1.Node{
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
		},
	}

	cluster := MapNodes(clusterName, nodes)

	assert.Equal(t, clusterName, cluster.ID)
	assert.Equal(t, clusterName, cluster.Data["clusterName"])
	assert.Equal(t, "cluster", cluster.Type)
}
