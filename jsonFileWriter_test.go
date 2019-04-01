package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSONFile(t *testing.T) {
	data := map[string][]map[string]interface{}{
		"ITComponent": []map[string]interface{}{
			{
				"name":                   "myapp",
				"uid":                    "a49ef9d4-5201-11e9-8647-d663bd873d93",
				"app.kubernetes.io/name": "myapp",
			},
			{
				"name":                   "otherapp",
				"uid":                    "bb53805a-5201-11e9-8647-d663bd873d93",
				"app.kubernetes.io/name": "otherapp",
			},
		},
	}
	filename := "writeJSONFileTestOutput.json"
	WriteJSONFile(data, filename)
	defer os.Remove(filename)
	outputFile, _ := ioutil.ReadFile(filename)
	var outputFileData map[string][]map[string]interface{}
	err := json.Unmarshal(outputFile, &outputFileData)
	if err != nil {
		t.Error(err)
	}
	assert.EqualValues(t, data, outputFileData)
}
