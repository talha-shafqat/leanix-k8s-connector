package mapper

import (
	"github.com/leanix/leanix-k8s-connector/pkg/set"

	corev1 "k8s.io/api/core/v1"
)

func aggregrateNodes(nodes *corev1.NodeList) map[string]interface{} {
	nodeAggregate := make(map[string]interface{})
	items := nodes.Items
	if len(items) == 0 {
		return nodeAggregate
	}
	availabilityZones := set.NewStringSet()
	nodeTypes := set.NewStringSet()

	for _, n := range items {
		availabilityZones.Add(n.Labels["failure-domain.beta.kubernetes.io/zone"])
		nodeTypes.Add(n.Labels["beta.kubernetes.io/instance-type"])
	}
	nodeAggregate["availabilityZones"] = availabilityZones.Items()
	nodeAggregate["dataCenter"] = items[0].Labels["failure-domain.beta.kubernetes.io/region"]
	nodeAggregate["nodeTypes"] = nodeTypes.Items()
	nodeAggregate["numberNodes"] = len(items)
	nodeAggregate["labels"] = labelSet(&items)
	return nodeAggregate
}

func labelSet(nodes *[]corev1.Node) map[string][]string {
	labelsAsSet := make(map[string]*set.String)
	labels := make(map[string][]string)
	for _, n := range *nodes {
		for l, v := range n.Labels {
			if labelsAsSet[l] == nil {
				labelsAsSet[l] = set.NewStringSet()
			}
			labelsAsSet[l].Add(v)
		}
	}
	for k, v := range labelsAsSet {
		labels[k] = v.Items()
	}
	return labels
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
