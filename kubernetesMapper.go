package main

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

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

// MapDeployments maps a kubernetes deployment list to a list of KubernetesObjects
func MapDeployments(deployments *appsv1.DeploymentList) []KubernetesObject {
	kubernetesObjects := make([]KubernetesObject, len(deployments.Items))
	for i, d := range deployments.Items {
		kubernetesObjects[i] = MapDeployment(d)
	}
	return kubernetesObjects
}

// MapDeployment maps a single kubernetes deployment to an KubernetesObject
func MapDeployment(deployment appsv1.Deployment) KubernetesObject {
	kubernetesObject := KubernetesObject{
		ID:   string(deployment.UID),
		Type: "deployment",
		Data: make(map[string]interface{}),
	}
	for k, v := range deployment.Labels {
		kubernetesObject.Data[k] = v
	}
	return kubernetesObject
}
