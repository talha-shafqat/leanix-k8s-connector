package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClusterKubernetesObject(t *testing.T) {
	clusterName := "mycluster"
	nodeInfo := KubernetesNodeInfo{
		DataCenter:        "westeurope",
		AvailabilityZones: []string{"0"},
		NodeTypes:         []string{"Standard_D4s_v3"},
		NumberNodes:       3,
		Labels:            map[string][]string{"agentpool": []string{"default"}},
	}

	cluster := NewClusterKubernetesObject(clusterName, nodeInfo)

	assert.Equal(t, clusterName, cluster.ID)
	assert.Equal(t, clusterName, cluster.Data["clusterName"])
	assert.Equal(t, "cluster", cluster.Type)
	assert.Equal(t, nodeInfo.DataCenter, cluster.Data["dataCenter"])
	assert.Equal(t, nodeInfo.AvailabilityZones, cluster.Data["availabilityZones"])
	assert.Equal(t, nodeInfo.NodeTypes, cluster.Data["nodeTypes"])
	assert.Equal(t, nodeInfo.NumberNodes, cluster.Data["numberNodes"])
	assert.Equal(t, nodeInfo.Labels["agentpool"], cluster.Data["labels"].(map[string][]string)["agentpool"])
}
