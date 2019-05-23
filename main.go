package main

import (
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/op/go-logging"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeConfigFlag       string = "kubeconfig"
	clusterNameFlag      string = "clustername"
	storageBackendFlag   string = "storage-backend"
	azureAccountNameFlag string = "azure-account-name"
	azureAccountKeyFlag  string = "azure-account-key"
	azureContainerFlag   string = "azure-container"
	verboseFlag          string = "verbose"
)

var log = logging.MustGetLogger("leanix-k8s-connector")

func main() {
	if home := homeDir(); home != "" {
		flag.String(kubeConfigFlag, filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.String(kubeConfigFlag, "", "absolute path to the kubeconfig file")
	}
	flag.String(clusterNameFlag, "", "unique name of the kubernets cluster")
	flag.String(storageBackendFlag, "file", "storage where the ldif.json file is placed. (file, azure)")
	flag.String(azureAccountNameFlag, "", "Azure storage account name")
	flag.String(azureAccountKeyFlag, "", "Azure storage account key")
	flag.String(azureContainerFlag, "", "Azure storage account container")
	flag.Bool(verboseFlag, false, "verbose log output")
	flag.Parse()
	// Let flags overwrite configs in viper
	viper.BindPFlags(flag.CommandLine)
	// Check for config values in env vars
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := InitLogger(viper.GetBool(verboseFlag))
	if err != nil {
		// use panic here because the logger functionality was not initalized
		panic(err)
	}
	log.Debugf("Target kubernetes cluster name: %s", viper.GetString(clusterNameFlag))

	if viper.GetString(clusterNameFlag) == "" {
		flag.PrintDefaults()
		log.Fatal("clustername flag must be set.")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString(kubeConfigFlag))
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Using kube config: %s", viper.GetString(kubeConfigFlag))
	log.Debugf("Kubernetes master from config: %s", config.Host)

	kubernetes, err := NewKubernetesAPI(config)
	if err != nil {
		log.Fatal(err)
	}
	blacklistedNamespaces := []string{"kube-system"}
	log.Debugf("Namespace blacklist: %v", blacklistedNamespaces)

	log.Debug("Get deployment list...")
	deployments, deploymentNodes, err := kubernetes.DeploymentsOnNodes(blacklistedNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Getting deployment list done.")

	log.Debug("Get statefulset list...")
	statefulsets, statefulsetNodes, err := kubernetes.StatefulSetsOnNodes(blacklistedNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Getting statefulset list done.")

	log.Debug("Listing nodes...")
	nodes, err := kubernetes.Nodes()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Listing nodes done.")

	log.Debug("Map nodes to kubernetes object")
	clusterKubernetesObject := NewClusterKubernetesObject(
		viper.GetString("clustername"),
		NewKubernetesNodeInfo(nodes),
	)

	log.Debug("Map deployments to kubernetes objects")
	deploymentKubernetesObjects := MapDeployments(viper.GetString(clusterNameFlag), deployments, deploymentNodes)
	log.Debug("Map statefulsets to kubernetes objects")
	statefulsetKubernetesObjects := MapStatefulSets(viper.GetString(clusterNameFlag), statefulsets, statefulsetNodes)

	kubernetesObjects := make([]KubernetesObject, 0)
	kubernetesObjects = append(kubernetesObjects, clusterKubernetesObject)
	kubernetesObjects = append(kubernetesObjects, deploymentKubernetesObjects...)
	kubernetesObjects = append(kubernetesObjects, statefulsetKubernetesObjects...)

	ldif := LDIF{
		ConnectorID:        "leanix-k8s-connector",
		ConnectorVersion:   "0.0.1",
		IntegrationVersion: "3",
		Description:        "Map kubernetes objects to LeanIX Fact Sheets",
		Content:            kubernetesObjects,
	}

	log.Debug("Marshal ldif")
	ldifByte, err := Marshal(ldif)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Upload ldif.json to %s", viper.GetString("storage-backend"))
	azureOpts := AzureStorageOpts{
		AccountName: viper.GetString(azureAccountNameFlag),
		AccountKey:  viper.GetString(azureAccountKeyFlag),
		Container:   viper.GetString(azureContainerFlag),
	}
	uploader, err := NewStorageBackend(viper.GetString("storage-backend"), &azureOpts)
	if err != nil {
		log.Fatal(err)
	}
	uploader.Upload(ldifByte)
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
