package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewKubernetesNodeInfo(t *testing.T) {
	// create a dummy nodes
	nodes := corev1.NodeList{
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
	nodeInfo := NewKubernetesNodeInfo(&nodes)

	assert.Equal(t, "westeurope", nodeInfo.DataCenter)
	assert.Len(t, nodeInfo.AvailabilityZones, 2)
	assert.Contains(t, nodeInfo.AvailabilityZones, "1")
	assert.Contains(t, nodeInfo.AvailabilityZones, "2")
	assert.Equal(t, 2, nodeInfo.NumberNodes)
	assert.Len(t, nodeInfo.NodeTypes, 2)
	assert.Contains(t, nodeInfo.NodeTypes, "Standard_D2s_v3")
	assert.Contains(t, nodeInfo.NodeTypes, "Standard_D8s_v3")
	// assert that Labels contains all labels present in any node object
	assert.Contains(t, nodeInfo.Labels["name"], "nodepool-1")
	assert.Contains(t, nodeInfo.Labels["name"], "nodepool-2")
	assert.Contains(t, nodeInfo.Labels["failure-domain.beta.kubernetes.io/region"], "westeurope")
	assert.Contains(t, nodeInfo.Labels["failure-domain.beta.kubernetes.io/zone"], "1")
	assert.Contains(t, nodeInfo.Labels["failure-domain.beta.kubernetes.io/zone"], "2")
	assert.Contains(t, nodeInfo.Labels["beta.kubernetes.io/instance-type"], "Standard_D8s_v3")
	assert.Contains(t, nodeInfo.Labels["beta.kubernetes.io/instance-type"], "Standard_D2s_v3")
}
