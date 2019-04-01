package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewOrchestrationFactSheet(t *testing.T) {
	clusterName := "awesome-k8s"
	nodeInfo := KubernetesNodeInfo{
		DataCenter:       "westeurope",
		AvailabilityZone: "0",
		NumberNodes:      "3",
		TypeNodes:        []string{"Standard_D2s_v3"},
	}
	factSheet := NewOrchestrationFactSheet(clusterName, nodeInfo)

	assert.Equal(t, factSheet["clusterName"], clusterName)
	assert.Equal(t, factSheet["type"], "Kubernetes")
	assert.Equal(t, factSheet["subFactSheetType"], "Orchestration")
	assert.Equal(t, factSheet["dataCenter"], "westeurope")
	assert.Equal(t, factSheet["availabilityZone"], "0")
	assert.Equal(t, factSheet["numberNodes"], "3")
	assert.Equal(t, factSheet["typeNodes"], []string{"Standard_D2s_v3"})
}

func TestNewFactSheet(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	deployment := NewDeployment("myapp", id, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})

	factSheet := NewFactSheet(deployment)

	assert.Equal(t, factSheet["name"], deployment.ObjectMeta.Name)
	assert.Equal(t, factSheet["uid"], deployment.ObjectMeta.UID)
	assert.Contains(t, factSheet, "app.kubernetes.io/name")
	assert.Equal(t, factSheet["app.kubernetes.io/name"], "myapp")
}

func TestGenerateFactSheets_Deployment_FactSheet(t *testing.T) {
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	otherAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	deployments := []apps.Deployment{
		NewDeployment("myapp", myAppID, map[string]string{
			"app.kubernetes.io/name": "myapp",
		}),
		NewDeployment("otherapp", otherAppID, map[string]string{
			"app.kubernetes.io/name": "otherapp",
		}),
	}

	factSheets := GenerateFactSheets(deployments)

	assert.Len(t, factSheets, 2)
}

func NewDeployment(name string, uuid uuid.UUID, labels map[string]string) apps.Deployment {
	uid := types.UID(uuid.String())
	return apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			UID:    uid,
			Labels: labels,
		},
	}
}
