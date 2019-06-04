package mapper

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAggregateNodes(t *testing.T) {
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
	expectedLabelAggregate := map[string][]string{
		"name": []string{"nodepool-1", "nodepool-2"},
		"failure-domain.beta.kubernetes.io/region": []string{"westeurope"},
		"failure-domain.beta.kubernetes.io/zone":   []string{"1", "2"},
		"beta.kubernetes.io/instance-type":         []string{"Standard_D2s_v3", "Standard_D8s_v3"},
	}

	nodeAggregate, err := aggregrateNodes(nodes)
	assert.NoError(t, err)

	assert.Equal(t, "westeurope", nodeAggregate["dataCenter"])
	assert.ElementsMatch(t, []string{"1", "2"}, nodeAggregate["availabilityZones"])
	assert.ElementsMatch(t, []string{"Standard_D2s_v3", "Standard_D8s_v3"}, nodeAggregate["nodeTypes"])
	assert.Equal(t, 2, nodeAggregate["numberNodes"])
	for k, v := range expectedLabelAggregate {
		assert.ElementsMatch(t, v, nodeAggregate["labels"].(map[string][]string)[k])
	}
}

func TestRedundant(t *testing.T) {
	type testExpected struct {
		multipleNodes bool
		zoneRedundant bool
	}
	tests := map[string]struct {
		input    []corev1.Node
		expected testExpected
	}{
		"single node": {
			input: []corev1.Node{
				corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "nodepool-1",
						Labels: map[string]string{
							"failure-domain.beta.kubernetes.io/zone": "1",
						},
					},
				},
			},
			expected: testExpected{
				multipleNodes: false,
				zoneRedundant: false,
			},
		},
		"multiple nodes": {
			input: []corev1.Node{
				corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "nodepool-1",
						Labels: map[string]string{
							"failure-domain.beta.kubernetes.io/zone": "1",
						},
					},
				},
				corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "nodepool-2",
						Labels: map[string]string{
							"failure-domain.beta.kubernetes.io/zone": "1",
						},
					},
				},
			},
			expected: testExpected{
				multipleNodes: true,
				zoneRedundant: false,
			},
		},
		"multiple zone redundant nodes": {
			input: []corev1.Node{
				corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "nodepool-1",
						Labels: map[string]string{
							"failure-domain.beta.kubernetes.io/zone": "1",
						},
					},
				},
				corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "nodepool-2",
						Labels: map[string]string{
							"failure-domain.beta.kubernetes.io/zone": "2",
						},
					},
				},
			},
			expected: testExpected{
				multipleNodes: true,
				zoneRedundant: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			multipleNodes, zoneRedundant := redundant(&test.input)
			assert.Equal(t, test.expected.multipleNodes, multipleNodes)
			assert.Equal(t, test.expected.zoneRedundant, zoneRedundant)
		})
	}

}

func TestAggregrateMemoryCapacity(t *testing.T) {
	oneGiB, err := resource.ParseQuantity("1Gi")
	if err != nil {
		t.Error(err)
	}
	fiveTwelveMiB, err := resource.ParseQuantity("512Mi")
	if err != nil {
		t.Error(err)
	}
	tests := map[string]struct {
		input    []corev1.Node
		expected float64
	}{
		"single 1Gi node": {
			input: []corev1.Node{
				corev1.Node{
					Status: corev1.NodeStatus{
						Capacity: corev1.ResourceList{
							corev1.ResourceMemory: oneGiB,
						},
					},
				},
			},
			expected: 1,
		},
		"two 1Gi nodes": {
			input: []corev1.Node{
				corev1.Node{
					Status: corev1.NodeStatus{
						Capacity: corev1.ResourceList{
							corev1.ResourceMemory: oneGiB,
						},
					},
				},
				corev1.Node{
					Status: corev1.NodeStatus{
						Capacity: corev1.ResourceList{
							corev1.ResourceMemory: oneGiB,
						},
					},
				},
			},
			expected: 2,
		},
		"512Mi and 1Gi nodes": {
			input: []corev1.Node{
				corev1.Node{
					Status: corev1.NodeStatus{
						Capacity: corev1.ResourceList{
							corev1.ResourceMemory: oneGiB,
						},
					},
				},
				corev1.Node{
					Status: corev1.NodeStatus{
						Capacity: corev1.ResourceList{
							corev1.ResourceMemory: fiveTwelveMiB,
						},
					},
				},
			},
			expected: 1.5,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			memory, err := aggregrateMemoryCapacity(&test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, memory)
		})
	}
}
