package main

import corev1 "k8s.io/api/core/v1"

// NewKubernetesNodeInfo creates a combined info struct from a list of nodes
func NewKubernetesNodeInfo(nodes *corev1.NodeList) KubernetesNodeInfo {
	items := nodes.Items
	if len(items) == 0 {
		return KubernetesNodeInfo{}
	}
	availabilityZones := NewStringSet()
	nodeTypes := NewStringSet()

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

// NewClusterKubernetesObject creates a KubernetesObject representing a cluster
func NewClusterKubernetesObject(clusterName string, nodeInfo KubernetesNodeInfo) KubernetesObject {
	return KubernetesObject{
		ID:   clusterName,
		Type: "cluster",
		Data: map[string]interface{}{
			"availabilityZones": nodeInfo.AvailabilityZones,
			"clusterName":       clusterName,
			"dataCenter":        nodeInfo.DataCenter,
			"nodeTypes":         nodeInfo.NodeTypes,
			"numberNodes":       nodeInfo.NumberNodes,
		},
	}
}
