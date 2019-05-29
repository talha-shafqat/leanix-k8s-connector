package mapper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

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

func TestMapDeployment(t *testing.T) {
	// init data which is constrant accross all tests
	clusterName := "mycluster"
	deploymentName := "myapp"
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}

	tests := map[string]struct {
		input    testInput
		expected testExpected
	}{
		"single replica": {
			input: testInput{
				id:                 id,
				replicas:           1,
				nodes:              1,
				nodesZoneRedundant: false,
			},
			expected: testExpected{
				isRedundant:                        false,
				isRedundantAcrossNodes:             false,
				isRedundantAcrossAvailabilityZones: false,
			},
		},
		"multiple replicas": {
			input: testInput{
				id:                 id,
				replicas:           2,
				nodes:              1,
				nodesZoneRedundant: false,
			},
			expected: testExpected{
				isRedundant:                        true,
				isRedundantAcrossNodes:             false,
				isRedundantAcrossAvailabilityZones: false,
			},
		},
		"multiple replicas on multiple nodes": {
			input: testInput{
				id:                 id,
				replicas:           2,
				nodes:              2,
				nodesZoneRedundant: false,
			},
			expected: testExpected{
				isRedundant:                        true,
				isRedundantAcrossNodes:             true,
				isRedundantAcrossAvailabilityZones: false,
			},
		},
		"multiple replicas on multiple zone redudant nodes": {
			input: testInput{
				id:                 id,
				replicas:           2,
				nodes:              2,
				nodesZoneRedundant: true,
			},
			expected: testExpected{
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
