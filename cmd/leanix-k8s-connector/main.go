package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leanix/leanix-k8s-connector/pkg/kubernetes"
	"github.com/leanix/leanix-k8s-connector/pkg/mapper"
	"github.com/leanix/leanix-k8s-connector/pkg/storage"
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
	localFilePathFlag    string = "local-file-path"
	verboseFlag          string = "verbose"
)

var log = logging.MustGetLogger("leanix-k8s-connector")

func main() {
	initLogger(viper.GetBool(verboseFlag))
	err := parseFlags()
	if err != nil {
		log.Critical(err)
	}
	log.Debugf("Target kubernetes cluster name: %s", viper.GetString(clusterNameFlag))

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString(kubeConfigFlag))
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Using kube config: %s", viper.GetString(kubeConfigFlag))
	log.Debugf("Kubernetes master from config: %s", config.Host)

	kubernetes, err := kubernetes.NewAPI(config)
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
	clusterKubernetesObject := mapper.MapNodes(
		viper.GetString("clustername"),
		nodes,
	)

	log.Debug("Map deployments to kubernetes objects")
	deploymentKubernetesObjects := mapper.MapDeployments(viper.GetString(clusterNameFlag), deployments, deploymentNodes)
	log.Debug("Map statefulsets to kubernetes objects")
	statefulsetKubernetesObjects := mapper.MapStatefulSets(viper.GetString(clusterNameFlag), statefulsets, statefulsetNodes)

	kubernetesObjects := make([]mapper.KubernetesObject, 0)
	kubernetesObjects = append(kubernetesObjects, clusterKubernetesObject)
	kubernetesObjects = append(kubernetesObjects, deploymentKubernetesObjects...)
	kubernetesObjects = append(kubernetesObjects, statefulsetKubernetesObjects...)

	ldif := mapper.LDIF{
		ConnectorID:        "leanix-k8s-connector",
		ConnectorVersion:   "0.0.1",
		IntegrationVersion: "3",
		Description:        "Map kubernetes objects to LeanIX Fact Sheets",
		Content:            kubernetesObjects,
	}

	log.Debug("Marshal ldif")
	ldifByte, err := storage.Marshal(ldif)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Upload ldif.json to %s", viper.GetString("storage-backend"))
	azureOpts := storage.AzureBlobOpts{
		AccountName: viper.GetString(azureAccountNameFlag),
		AccountKey:  viper.GetString(azureAccountKeyFlag),
		Container:   viper.GetString(azureContainerFlag),
	}
	localFileOpts := storage.LocalFileOpts{
		Path: viper.GetString(localFilePathFlag),
	}
	uploader, err := storage.NewBackend(viper.GetString("storage-backend"), &azureOpts, &localFileOpts)
	if err != nil {
		log.Fatal(err)
	}
	uploader.Upload(ldifByte)
}

func parseFlags() error {
	if home := homeDir(); home != "" {
		flag.String(kubeConfigFlag, filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.String(kubeConfigFlag, "", "absolute path to the kubeconfig file")
	}
	flag.String(clusterNameFlag, "", "unique name of the kubernets cluster")
	flag.String(storageBackendFlag, storage.FileStorage, fmt.Sprintf("storage where the ldif.json file is placed. (%s, %s)", storage.FileStorage, storage.AzureBlobStorage))
	flag.String(azureAccountNameFlag, "", "Azure storage account name")
	flag.String(azureAccountKeyFlag, "", "Azure storage account key")
	flag.String(azureContainerFlag, "", "Azure storage account container")
	flag.String(localFilePathFlag, ".", "path to place the ldif file when using local file storage backend")
	flag.Bool(verboseFlag, false, "verbose log output")
	flag.Parse()
	// Let flags overwrite configs in viper
	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		return err
	}
	// Check for config values in env vars
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	if viper.GetString(clusterNameFlag) == "" {
		return errors.New("clustername flag must be set")
	}
	return nil
}

// InitLogger initialise the logger for stdout and log file
func initLogger(verbose bool) {
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
		log.Warningf("unable to log to 'leanix-k8s-connector.log': %s\n", err)
	}
	fileLogger := logging.NewLogBackend(f, "", 0)
	logging.SetBackend(fileLogger, stdoutLeveled)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
