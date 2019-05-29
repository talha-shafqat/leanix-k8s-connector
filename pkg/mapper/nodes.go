package mapper

import (
	"github.com/leanix/leanix-k8s-connector/pkg/set"

	corev1 "k8s.io/api/core/v1"
)

// KubernetesNodeInfo holds meta information about a kubernetes cluster
type KubernetesNodeInfo struct {
	DataCenter        string
	AvailabilityZones []string
	NumberNodes       int
	NodeTypes         []string
}

// NewKubernetesNodeInfo creates a combined info struct from a list of nodes
func NewKubernetesNodeInfo(nodes *corev1.NodeList) KubernetesNodeInfo {
	items := nodes.Items
	if len(items) == 0 {
		return KubernetesNodeInfo{}
	}
	availabilityZones := set.NewStringSet()
	nodeTypes := set.NewStringSet()

	for _, n := range items {
		availabilityZones.Add(n.Labels["failure-domain.beta.kubernetes.io/zone"])
		nodeTypes.Add(n.Labels["beta.kubernetes.io/instance-type"])
	}

	return KubernetesNodeInfo{
		DataCenter:        items[0].Labels["failure-domain.beta.kubernetes.io/region"],
		AvailabilityZones: availabilityZones.Items(),
		NumberNodes:       len(items),
		NodeTypes:         nodeTypes.Items(),
	}
}

func redundant(nodes *[]corev1.Node) (bool, bool) {
	nodeNames := set.NewStringSet()
	availabilityZones := set.NewStringSet()
	for _, n := range *nodes {
		nodeNames.Add(n.GetName())
		availabilityZones.Add(n.Labels["failure-domain.beta.kubernetes.io/zone"])
	}
	return len(nodeNames.Items()) > 1, len(availabilityZones.Items()) > 1
}
