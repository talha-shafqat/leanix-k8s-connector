package mapper

import (
	corev1 "k8s.io/api/core/v1"
)

// MapNodes mapps a list of nodes and a given cluster name into a KubernetesObject.
// In the process it aggregates the information from muliple nodes into one cluster object.
func MapNodes(clusterName string, nodes *corev1.NodeList) (*KubernetesObject, error) {
	nodeAggregate, err := aggregrateNodes(nodes)
	if err != nil {
		return nil, err
	}
	nodeAggregate["clusterName"] = clusterName
	return &KubernetesObject{
		ID:   clusterName,
		Type: "cluster",
		Data: nodeAggregate,
	}, nil
}
