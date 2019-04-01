package main

import appsv1 "k8s.io/api/apps/v1"

// NewOrchestrationFactSheet creates a new orchestration FactSheet, which
// represents the kubernetes cluster.
func NewOrchestrationFactSheet(clusterName string, nodeInfo *KubernetesNodeInfo) FactSheet {
	fs := make(FactSheet)
	fs["clusterName"] = clusterName
	fs["type"] = "Kubernetes"
	fs["subFactSheetType"] = "Orchestration"
	fs["dataCenter"] = nodeInfo.DataCenter
	fs["availabilityZones"] = nodeInfo.AvailabilityZones
	fs["numberNodes"] = nodeInfo.NumberNodes
	fs["nodeTypes"] = nodeInfo.NodeTypes
	return fs
}

// NewFactSheet creates a new Fact Sheet from a deployment
func NewFactSheet(d appsv1.Deployment) FactSheet {
	fs := make(FactSheet)
	fs["uid"] = d.UID
	fs["name"] = d.Name
	for labelName, labelValue := range d.Labels {
		fs[labelName] = labelValue
	}
	return fs
}

// GenerateFactSheets takes a list of deployments extracts the uid, name and labels into
// a map.
func GenerateFactSheets(deployments []appsv1.Deployment) []FactSheet {
	factSheets := make([]FactSheet, 0)
	for _, d := range deployments {
		factSheets = append(factSheets, NewFactSheet(d))
	}
	return factSheets
}
