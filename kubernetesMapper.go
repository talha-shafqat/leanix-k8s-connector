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

// MapStatefulSets maps a kubernetes statefulset list to a list of KubernetesObjects
func MapStatefulSets(clusterName string, statefulsets *appsv1.StatefulSetList, nodes *[]corev1.Node) []KubernetesObject {
	kubernetesObjects := make([]KubernetesObject, len(statefulsets.Items))
	for i, s := range statefulsets.Items {
		kubernetesObjects[i] = MapStatefulSet(clusterName, s, nodes)
	}
	return kubernetesObjects
}

// MapDeployments maps a kubernetes deployment list to a list of KubernetesObjects
func MapDeployments(clusterName string, deployments *appsv1.DeploymentList, nodes *[]corev1.Node) []KubernetesObject {
	kubernetesObjects := make([]KubernetesObject, len(deployments.Items))
	for i, d := range deployments.Items {
		kubernetesObjects[i] = MapDeployment(clusterName, d, nodes)
	}
	return kubernetesObjects
}

// MapDeployment maps a single kubernetes deployment to an KubernetesObject
func MapDeployment(clusterName string, deployment appsv1.Deployment, nodes *[]corev1.Node) KubernetesObject {
	kubernetesObject := KubernetesObject{
		ID:   string(deployment.UID),
		Type: "deployment",
		Data: make(map[string]interface{}),
	}
	for k, v := range deployment.Labels {
		kubernetesObject.Data[k] = v
	}
	redundantAcrossNodes, redundantAcrossAvailabilityZones := redundant(nodes)
	kubernetesObject.Data["clusterName"] = clusterName
	kubernetesObject.Data["isStateful"] = false
	kubernetesObject.Data["isRedundant"] = deployment.Status.Replicas > 1
	kubernetesObject.Data["isRedundantAcrossNodes"] = redundantAcrossNodes
	kubernetesObject.Data["isRedundantAcrossAvailabilityZones"] = redundantAcrossAvailabilityZones
	return kubernetesObject
}

// MapStatefulSet maps a single kubernetes StatefulSet to an KubernetesObject
func MapStatefulSet(clusterName string, statefulset appsv1.StatefulSet, nodes *[]corev1.Node) KubernetesObject {
	kubernetesObject := KubernetesObject{
		ID:   string(statefulset.UID),
		Type: "statefulSet",
		Data: make(map[string]interface{}),
	}
	for k, v := range statefulset.Labels {
		kubernetesObject.Data[k] = v
	}
	redundantAcrossNodes, redundantAcrossAvailabilityZones := redundant(nodes)
	kubernetesObject.Data["clusterName"] = clusterName
	kubernetesObject.Data["isStateful"] = true
	kubernetesObject.Data["isRedundant"] = statefulset.Status.Replicas > 1
	kubernetesObject.Data["isRedundantAcrossNodes"] = redundantAcrossNodes
	kubernetesObject.Data["isRedundantAcrossAvailabilityZones"] = redundantAcrossAvailabilityZones
	return kubernetesObject
}

func redundant(nodes *[]corev1.Node) (bool, bool) {
	nodeNames := NewStringSet()
	availabilityZones := NewStringSet()
	for _, n := range *nodes {
		nodeNames.Add(n.GetName())
		availabilityZones.Add(n.Labels["failure-domain.beta.kubernetes.io/zone"])
	}
	return len(nodeNames.Items()) > 1, len(availabilityZones.Items()) > 1
}
