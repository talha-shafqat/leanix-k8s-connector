package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSONFile(t *testing.T) {
	data := LDIF{
		ConnectorID:        "leanix-k8s-connector",
		ConnectorVersion:   "0.0.1",
		IntegrationVersion: "3",
		Description:        "Map kubernetes objects to LeanIX Fact Sheets",
		Content: []KubernetesObject{
			KubernetesObject{
				ID:   "my-cluster",
				Type: "cluster",
				Data: map[string]interface{}{
					"availabilityZones": []interface{}{"0"},
					"clusterName":       "my-cluster",
					"dataCenter":        "westeurope",
					"nodeTypes":         []interface{}{"Standard_D4s_v3"},
					// Cast to float64 because unmarshaling uses float64 when interface{} is used
					"numberNodes": float64(3),
				},
			},
			KubernetesObject{
				ID:   "a49ef9d4-5201-11e9-8647-d663bd873d93",
				Type: "deployment",
				Data: map[string]interface{}{
					"app.kubernetes.io/name": "myapp",
				},
			},
		},
	}
	filename := "writeJSONFileTestOutput.json"
	WriteJSONFile(data, filename)
	defer os.Remove(filename)
	outputFile, _ := ioutil.ReadFile(filename)
	var ldif LDIF
	err := json.Unmarshal(outputFile, &ldif)
	if err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, data, ldif)
}
