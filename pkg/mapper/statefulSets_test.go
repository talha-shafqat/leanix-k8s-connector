package mapper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

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
	// init data which is constrant accross all tests
	clusterName := "mycluster"
	name := "myapp"
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
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			deployment := newStatefulSet(name, test.input.id, test.input.replicas, map[string]string{
				"app.kubernetes.io/name": name,
			})
			nodes := newNodes(test.input.nodes, test.input.nodesZoneRedundant)
			result := MapStatefulSet(clusterName, deployment, nodes)

			assert.Equal(t, result.ID, id.String())
			assert.Equal(t, result.Type, "statefulSet")
			assert.Equal(t, result.Data["app_kubernetes_io_name"], name)
			assert.Equal(t, result.Data["clusterName"], clusterName)
			assert.Equal(t, result.Data["isStateful"], true)
			assert.Equal(t, result.Data["isRedundant"], test.expected.isRedundant)
			assert.Equal(t, result.Data["isRedundantAcrossNodes"], test.expected.isRedundantAcrossNodes)
			assert.Equal(t, result.Data["isRedundantAcrossAvailabilityZones"], test.expected.isRedundantAcrossAvailabilityZones)
		})
	}
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
