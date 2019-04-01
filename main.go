package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	var clusterName *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	clusterName = flag.String("clustername", "", "unique name of the kubernets cluster")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	kubernetes, err := NewKubernetesAPI(config)
	if err != nil {
		log.Fatal(err)
	}
	deployments, err := kubernetes.Deployments([]string{"kube-system"})
	if err != nil {
		log.Fatal(err)
	}
	nodes, err := kubernetes.Nodes()
	if err != nil {
		log.Fatal(err)
	}

	orchestrationFactSheet := NewOrchestrationFactSheet(
		*clusterName,
		NewKubernetesNodeInfo(nodes),
	)
	factSheets := []FactSheet{orchestrationFactSheet}
	deploymentUIDs := make([]interface{}, 0)
	for _, fs := range GenerateFactSheets(deployments.Items) {
		factSheets = append(factSheets, fs)
		deploymentUIDs = append(deploymentUIDs, fs["uid"])
	}

	factSheetJSON := map[string]interface{}{
		"ITComponent": factSheets,
	}
	err = WriteJSONFile(factSheetJSON, "factSheets.json")
	if err != nil {
		log.Fatal(err)
	}
	relationJSON := Relations(orchestrationFactSheet["clusterName"].(string), deploymentUIDs)
	err = WriteJSONFile(relationJSON, "relations.json")
	if err != nil {
		log.Fatal(err)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
