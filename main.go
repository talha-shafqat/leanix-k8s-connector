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
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
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

	// orchestrationFactSheet := NewOrchestrationFactSheet()
	factSheets := GenerateFactSheets(deployments.Items)
	factSheetJSON := map[string]interface{}{
		"ITComponent": factSheets,
	}
	err = WriteJSONFile(factSheetJSON, "factSheets.json")
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
