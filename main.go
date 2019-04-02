package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/op/go-logging"
	"k8s.io/client-go/tools/clientcmd"
)

var log = logging.MustGetLogger("leanix-k8s-connector")

func main() {
	// Parse flags
	var kubeconfig *string
	var clusterName *string
	var verbose *bool
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	clusterName = flag.String("clustername", "", "unique name of the kubernets cluster")
	verbose = flag.Bool("verbose", false, "Verbose log output.")
	flag.Parse()
	err := InitLogger(*verbose)
	if err != nil {
		// use panic here because the logger functionality was not initalized
		panic(err)
	}
	log.Debugf("Target kubernetes cluster name: %s", *clusterName)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Using kube config: %s", *kubeconfig)
	log.Debugf("Kubernetes master from config: %s", config.Host)

	kubernetes, err := NewKubernetesAPI(config)
	if err != nil {
		log.Fatal(err)
	}
	blacklist := []string{"kube-system"}
	log.Debugf("Namespace blacklist: %v", blacklist)
	log.Debug("Get deployment list...")
	deployments, err := kubernetes.Deployments(blacklist)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Getting deployment list done.")
	log.Debug("Listing nodes...")
	nodes, err := kubernetes.Nodes()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Listing nodes done.")

	log.Debug("Building orchestration Fact Sheet from kubernetes nodes...")
	orchestrationFactSheet := NewOrchestrationFactSheet(
		*clusterName,
		NewKubernetesNodeInfo(nodes),
	)
	log.Debug("Building orchestration Fact Sheet done.")
	factSheets := []FactSheet{orchestrationFactSheet}
	deploymentUIDs := make([]interface{}, 0)
	log.Debug("Generating Fact Sheets from kubernetes deployments...")
	for _, fs := range GenerateFactSheets(deployments.Items) {
		factSheets = append(factSheets, fs)
		deploymentUIDs = append(deploymentUIDs, fs["uid"])
	}
	log.Debug("Generating Fact Sheets from kubernetes deployments done.")

	factSheetJSON := map[string]interface{}{
		"ITComponent": factSheets,
	}
	log.Debug("Write factSheets.json to disk...")
	err = WriteJSONFile(factSheetJSON, "factSheets.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Write factSheets.json to disk done.")
	log.Debug("Write relations.json to disk...")
	relationJSON := Relations(orchestrationFactSheet["clusterName"].(string), deploymentUIDs)
	err = WriteJSONFile(relationJSON, "relations.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Write relations.json to disk done.")
}

// InitLogger initialise the logger for stdout and log file
func InitLogger(verbose bool) error {
	format := logging.MustStringFormatter(`%{time:15:04:05.000} ▶ [%{level:.4s}] %{message}`)
	logging.SetFormatter(format)

	// stdout logging backend
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutLeveled := logging.AddModuleLevel(stdout)
	if verbose {
		stdoutLeveled.SetLevel(logging.DEBUG, "")
	} else {
		stdoutLeveled.SetLevel(logging.INFO, "")
	}

	// file logging backend
	f, err := os.OpenFile("leanix-k8s-connector.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	fileLogger := logging.NewLogBackend(f, "", 0)
	logging.SetBackend(fileLogger, stdoutLeveled)

	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
