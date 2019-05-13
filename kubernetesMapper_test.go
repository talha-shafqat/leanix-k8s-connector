package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestNewKubernetesNodeInfo(t *testing.T) {
	// create a dummy nodes
	nodes := corev1.NodeList{
		Items: []corev1.Node{
			corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nodepool-1",
					Labels: map[string]string{
						"name": "nodepool-1",
						"failure-domain.beta.kubernetes.io/region": "westeurope",
						"failure-domain.beta.kubernetes.io/zone":   "1",
						"beta.kubernetes.io/instance-type":         "Standard_D2s_v3",
					},
				},
			},
			corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nodepool-2",
					Labels: map[string]string{
						"name": "nodepool-2",
						"failure-domain.beta.kubernetes.io/region": "westeurope",
						"failure-domain.beta.kubernetes.io/zone":   "2",
						"beta.kubernetes.io/instance-type":         "Standard_D8s_v3",
					},
				},
			},
		},
	}
	nodeInfo := NewKubernetesNodeInfo(&nodes)

	assert.Equal(t, "westeurope", nodeInfo.DataCenter)
	assert.Len(t, nodeInfo.AvailabilityZones, 2)
	assert.Contains(t, nodeInfo.AvailabilityZones, "1")
	assert.Contains(t, nodeInfo.AvailabilityZones, "2")
	assert.Equal(t, 2, nodeInfo.NumberNodes)
	assert.Len(t, nodeInfo.NodeTypes, 2)
	assert.Contains(t, nodeInfo.NodeTypes, "Standard_D2s_v3")
	assert.Contains(t, nodeInfo.NodeTypes, "Standard_D8s_v3")
}

func TestNewClusterKubernetesObject(t *testing.T) {
	clusterName := "mycluster"
	nodeInfo := KubernetesNodeInfo{
		DataCenter:        "westeurope",
		AvailabilityZones: []string{"0"},
		NodeTypes:         []string{"Standard_D4s_v3"},
		NumberNodes:       3,
	}

	cluster := NewClusterKubernetesObject(clusterName, nodeInfo)

	assert.Equal(t, clusterName, cluster.ID)
	assert.Equal(t, clusterName, cluster.Data["clusterName"])
	assert.Equal(t, "cluster", cluster.Type)
	assert.Equal(t, nodeInfo.DataCenter, cluster.Data["dataCenter"])
	assert.Equal(t, nodeInfo.AvailabilityZones, cluster.Data["availabilityZones"])
	assert.Equal(t, nodeInfo.NodeTypes, cluster.Data["nodeTypes"])
	assert.Equal(t, nodeInfo.NumberNodes, cluster.Data["numberNodes"])
}

func TestMapDeployments(t *testing.T) {
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	otherAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	deployments := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			newDeployment("myapp", myAppID, 1, map[string]string{
				"app.kubernetes.io/name": "myapp",
			}),
			newDeployment("otherapp", otherAppID, 2, map[string]string{
				"app.kubernetes.io/name": "otherapp",
			}),
		},
	}

	ko := MapDeployments(deployments)

	assert.Len(t, ko, 2)
	assert.Equal(t, ko[0].ID, myAppID.String())
	assert.Equal(t, ko[1].ID, otherAppID.String())
}

func TestMapDeployment_singleReplica(t *testing.T) {
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	deployment := newDeployment("myapp", myAppID, 1, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})

	ko := MapDeployment(deployment)

	assert.Equal(t, ko.ID, myAppID.String())
	assert.Equal(t, ko.Type, "deployment")
	assert.Equal(t, ko.Data["app.kubernetes.io/name"], "myapp")
	assert.Equal(t, ko.Data["isRedundant"], false)
}

func TestMapDeployment_mutlipleReplicas(t *testing.T) {
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	deployment := newDeployment("myapp", myAppID, 2, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})

	ko := MapDeployment(deployment)

	assert.Equal(t, ko.ID, myAppID.String())
	assert.Equal(t, ko.Type, "deployment")
	assert.Equal(t, ko.Data["app.kubernetes.io/name"], "myapp")
	assert.Equal(t, ko.Data["isRedundant"], true)
}

func newDeployment(name string, uuid uuid.UUID, replicas int32, labels map[string]string) appsv1.Deployment {
	uid := types.UID(uuid.String())
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			UID:    uid,
			Labels: labels,
		},
		Status: appsv1.DeploymentStatus{
			Replicas: replicas,
		},
	}
}
