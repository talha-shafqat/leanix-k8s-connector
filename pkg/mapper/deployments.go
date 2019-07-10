package mapper

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

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
		kubernetesObject.Data[replacer.Replace(k)] = v
	}
	redundantAcrossNodes, redundantAcrossAvailabilityZones := redundant(nodes)
	kubernetesObject.Data["clusterName"] = clusterName
	kubernetesObject.Data["isStateful"] = false
	kubernetesObject.Data["isRedundant"] = deployment.Status.Replicas > 1
	kubernetesObject.Data["isRedundantAcrossNodes"] = redundantAcrossNodes
	kubernetesObject.Data["isRedundantAcrossAvailabilityZones"] = redundantAcrossAvailabilityZones
	return kubernetesObject
}
