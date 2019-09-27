package mapper

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// MapStatefulSets maps a kubernetes statefulset list to a list of KubernetesObjects
func MapStatefulSets(clusterName string, statefulsets *appsv1.StatefulSetList, nodes *[]corev1.Node) []KubernetesObject {
	kubernetesObjects := make([]KubernetesObject, len(statefulsets.Items))
	for i, s := range statefulsets.Items {
		kubernetesObjects[i] = MapStatefulSet(clusterName, s, nodes)
	}
	return kubernetesObjects
}

// MapStatefulSet maps a single Kubernetes StatefulSet to an KubernetesObject
func MapStatefulSet(clusterName string, statefulset appsv1.StatefulSet, nodes *[]corev1.Node) KubernetesObject {
	kubernetesObject := KubernetesObject{
		ID:   string(statefulset.UID),
		Type: "statefulSet",
		Data: make(map[string]interface{}),
	}
	for k, v := range statefulset.Labels {
		kubernetesObject.Data[replacer.Replace(k)] = v
	}
	redundantAcrossNodes, redundantAcrossAvailabilityZones := redundant(nodes)
	kubernetesObject.Data["clusterName"] = clusterName
	kubernetesObject.Data["isStateful"] = true
	kubernetesObject.Data["isRedundant"] = statefulset.Status.Replicas > 1
	kubernetesObject.Data["isRedundantAcrossNodes"] = redundantAcrossNodes
	kubernetesObject.Data["isRedundantAcrossAvailabilityZones"] = redundantAcrossAvailabilityZones
	return kubernetesObject
}
