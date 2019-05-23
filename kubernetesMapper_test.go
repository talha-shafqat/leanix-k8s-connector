package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type NameAndZone struct {
	Name string
	Zone string
}

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
	clusterName := "mycluster"
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
	node := newNodes(1, false)

	ko := MapDeployments(clusterName, deployments, node)

	assert.Len(t, ko, 2)
	assert.Equal(t, ko[0].ID, myAppID.String())
	assert.Equal(t, ko[1].ID, otherAppID.String())
}

type deploymentMapTestInput struct {
	id                 uuid.UUID
	replicas           int32
	nodes              int
	nodesZoneRedundant bool
}

type deploymentMapTestExpected struct {
	isRedundant                        bool
	isRedundantAcrossNodes             bool
	isRedundantAcrossAvailabilityZones bool
}

func TestMapDeployment(t *testing.T) {
	// init data which is constrant accross all tests
	clusterName := "mycluster"
	deploymentName := "myapp"
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}

	tests := map[string]struct {
		input    deploymentMapTestInput
		expected deploymentMapTestExpected
	}{
		"single replica": {
			input: deploymentMapTestInput{
				id:                 id,
				replicas:           1,
				nodes:              1,
				nodesZoneRedundant: false,
			},
			expected: deploymentMapTestExpected{
				isRedundant:                        false,
				isRedundantAcrossNodes:             false,
				isRedundantAcrossAvailabilityZones: false,
			},
		},
		"multiple replicas": {
			input: deploymentMapTestInput{
				id:                 id,
				replicas:           2,
				nodes:              1,
				nodesZoneRedundant: false,
			},
			expected: deploymentMapTestExpected{
				isRedundant:                        true,
				isRedundantAcrossNodes:             false,
				isRedundantAcrossAvailabilityZones: false,
			},
		},
		"multiple replicas on multiple nodes": {
			input: deploymentMapTestInput{
				id:                 id,
				replicas:           2,
				nodes:              2,
				nodesZoneRedundant: false,
			},
			expected: deploymentMapTestExpected{
				isRedundant:                        true,
				isRedundantAcrossNodes:             true,
				isRedundantAcrossAvailabilityZones: false,
			},
		},
		"multiple replicas on multiple zone redudant nodes": {
			input: deploymentMapTestInput{
				id:                 id,
				replicas:           2,
				nodes:              2,
				nodesZoneRedundant: true,
			},
			expected: deploymentMapTestExpected{
				isRedundant:                        true,
				isRedundantAcrossNodes:             true,
				isRedundantAcrossAvailabilityZones: true,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			deployment := newDeployment(deploymentName, test.input.id, test.input.replicas, map[string]string{
				"app.kubernetes.io/name": deploymentName,
			})
			nodes := newNodes(test.input.nodes, test.input.nodesZoneRedundant)
			result := MapDeployment(clusterName, deployment, nodes)

			assert.Equal(t, result.ID, id.String())
			assert.Equal(t, result.Type, "deployment")
			assert.Equal(t, result.Data["app.kubernetes.io/name"], deploymentName)
			assert.Equal(t, result.Data["clusterName"], clusterName)
			assert.Equal(t, result.Data["isStateful"], false)
			assert.Equal(t, result.Data["isRedundant"], test.expected.isRedundant)
			assert.Equal(t, result.Data["isRedundantAcrossNodes"], test.expected.isRedundantAcrossNodes)
			assert.Equal(t, result.Data["isRedundantAcrossAvailabilityZones"], test.expected.isRedundantAcrossAvailabilityZones)
		})
	}
}

func TestMapStatefulSets(t *testing.T) {
	clusterName := "mycluster"
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	otherAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	statefulsets := &appsv1.StatefulSetList{
		Items: []appsv1.StatefulSet{
			newStatefulSet("myapp", myAppID, 1, map[string]string{
				"app.kubernetes.io/name": "myapp",
			}),
			newStatefulSet("otherapp", otherAppID, 1, map[string]string{
				"app.kubernetes.io/name": "otherapp",
			}),
		},
	}
	node := newNodes(1, false)

	ko := MapStatefulSets(clusterName, statefulsets, node)

	assert.Len(t, ko, 2)
	assert.Equal(t, ko[0].ID, myAppID.String())
	assert.Equal(t, ko[1].ID, otherAppID.String())
}

func TestMapStatefulSet(t *testing.T) {
	clusterName := "mycluster"
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	statefulset := newStatefulSet("myapp", myAppID, 1, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})
	node := newNodes(1, false)

	ko := MapStatefulSet(clusterName, statefulset, node)

	assert.Equal(t, ko.ID, myAppID.String())
	assert.Equal(t, ko.Type, "statefulSet")
	assert.Equal(t, ko.Data["isStateful"], true)
	assert.Equal(t, ko.Data["app.kubernetes.io/name"], "myapp")
	assert.Equal(t, ko.Data["clusterName"], clusterName)
	assert.Equal(t, ko.Data["isRedundant"], false)
	assert.Equal(t, ko.Data["isRedundantAcrossNodes"], false)
	assert.Equal(t, ko.Data["isRedundantAcrossAvailabilityZones"], false)
}

func TestMapStatefulSet_multipleReplicas(t *testing.T) {
	clusterName := "mycluster"
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	statefulset := newStatefulSet("myapp", myAppID, 2, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})
	node := newNodes(1, false)

	ko := MapStatefulSet(clusterName, statefulset, node)

	assert.Equal(t, ko.ID, myAppID.String())
	assert.Equal(t, ko.Type, "statefulSet")
	assert.Equal(t, ko.Data["isStateful"], true)
	assert.Equal(t, ko.Data["app.kubernetes.io/name"], "myapp")
	assert.Equal(t, ko.Data["clusterName"], clusterName)
	assert.Equal(t, ko.Data["isRedundant"], true)
	assert.Equal(t, ko.Data["isRedundantAcrossNodes"], false)
	assert.Equal(t, ko.Data["isRedundantAcrossAvailabilityZones"], false)
}

func TestMapStatefulSet_multipleReplicas_multipleNodes(t *testing.T) {
	clusterName := "mycluster"
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	statefulset := newStatefulSet("myapp", myAppID, 2, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})
	node := newNodes(2, false)

	ko := MapStatefulSet(clusterName, statefulset, node)

	assert.Equal(t, ko.ID, myAppID.String())
	assert.Equal(t, ko.Type, "statefulSet")
	assert.Equal(t, ko.Data["isStateful"], true)
	assert.Equal(t, ko.Data["app.kubernetes.io/name"], "myapp")
	assert.Equal(t, ko.Data["clusterName"], clusterName)
	assert.Equal(t, ko.Data["isRedundant"], true)
	assert.Equal(t, ko.Data["isRedundantAcrossNodes"], true)
	assert.Equal(t, ko.Data["isRedundantAcrossAvailabilityZones"], false)
}

func TestMapStatefulSet_multipleReplicas_multipleNodes_multipleRegions(t *testing.T) {
	clusterName := "mycluster"
	myAppID, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	statefulset := newStatefulSet("myapp", myAppID, 2, map[string]string{
		"app.kubernetes.io/name": "myapp",
	})
	node := newNodes(2, true)

	ko := MapStatefulSet(clusterName, statefulset, node)

	assert.Equal(t, ko.ID, myAppID.String())
	assert.Equal(t, ko.Type, "statefulSet")
	assert.Equal(t, ko.Data["isStateful"], true)
	assert.Equal(t, ko.Data["app.kubernetes.io/name"], "myapp")
	assert.Equal(t, ko.Data["clusterName"], clusterName)
	assert.Equal(t, ko.Data["isRedundant"], true)
	assert.Equal(t, ko.Data["isRedundantAcrossNodes"], true)
	assert.Equal(t, ko.Data["isRedundantAcrossAvailabilityZones"], true)
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

func newNodes(numberOfNodes int, nodesZoneRedundant bool) *[]corev1.Node {
	nodes := make([]corev1.Node, numberOfNodes)
	for i := range nodes {
		name := fmt.Sprintf("kubelet-%d", i)
		zone := "0"
		if nodesZoneRedundant {
			zone = string(i)
		}
		nodes[i] = corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name":                                   name,
					"failure-domain.beta.kubernetes.io/zone": zone,
				},
			}}
	}
	return &nodes
}

func newStatefulSet(name string, uuid uuid.UUID, replicas int32, labels map[string]string) appsv1.StatefulSet {
	uid := types.UID(uuid.String())
	return appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			UID:    uid,
			Labels: labels,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas: replicas,
		},
	}
}
