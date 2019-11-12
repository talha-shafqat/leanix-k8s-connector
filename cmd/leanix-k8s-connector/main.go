package main

import (
	"bytes"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	clusterNameFlag         string = "clustername"
	storageBackendFlag      string = "storage-backend"
	azureAccountNameFlag    string = "azure-account-name"
	azureAccountKeyFlag     string = "azure-account-key"
	azureContainerFlag      string = "azure-container"
	localFilePathFlag       string = "local-file-path"
	verboseFlag             string = "verbose"
	connectorIDFlag         string = "connector-id"
	blacklistNamespacesFlag string = "blacklist-namespaces"
	lxWorkspaceFlag         string = "lx-workspace"
	localFlag               string = "local"
)

const connectorVersion string = "1.1.0"
const lxVersion string = "1.0.0"

var log = logging.MustGetLogger("leanix-k8s-connector")

func main() {
	stdoutLogger, debugLogBuffer := initLogger()
	err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	enableVerbose(stdoutLogger, viper.GetBool(verboseFlag))
	log.Info("----------Start----------")
	log.Infof("LeanIX connector version: %s", connectorVersion)
	log.Infof("LeanIX integration version: %s", lxVersion)
	log.Infof("Target LeanIX workspace: %s", viper.GetString(lxWorkspaceFlag))
	log.Infof("Target Kubernetes cluster name: %s", viper.GetString(clusterNameFlag))

	var config *restclient.Config
	if viper.GetBool(localFlag) {
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// use the current context in kubeconfig
		config, err = restclient.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to load kube config. Running in Kubernetes?\n%s", err)
		}
	}

	log.Debugf("Kubernetes master from config: %s", config.Host)

	kubernetes, err := kubernetes.NewAPI(config)
	if err != nil {
		log.Fatal(err)
	}
	dynClient, err := dynamic.NewForConfig(config)

	resourcesList, err := ServerPreferredListableResources(kubernetes.Client.Discovery())
	if err != nil {
		log.Fatal(err)
	}
	groupVersionResources, err := discovery.GroupVersionResources(resourcesList)
	if err != nil {
		log.Panic(err)
	}

	kubernetesObjects := make([]mapper.KubernetesObject, 0)
	for gvr := range groupVersionResources {
		if gvr.Resource == "secret" {
			continue
		}
		instances, err := dynClient.Resource(gvr).List(metav1.ListOptions{})
		if err != nil {
			log.Panic(err)
		}
		for _, i := range instances.Items {
			nko := mapper.KubernetesObject{
				Type: i.GetKind(),
				ID:   string(i.GetUID()),
				Data: i.Object,
			}
			kubernetesObjects = append(kubernetesObjects, nko)
		}
		log.Info(instances)
	}

	// log.Debug("Get blacklist namespaces list...")
	// blacklistedNamespacesList := viper.GetStringSlice(blacklistNamespacesFlag)
	// blacklistedNamespaces, err := kubernetes.Namespaces(blacklistedNamespacesList)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Debug("Getting blacklist namespaces list done.")
	// log.Infof("Namespace blacklist: %v", blacklistedNamespaces)

	// log.Debug("Get deployment list...")
	// deployments, deploymentNodes, err := kubernetes.DeploymentsOnNodes(blacklistedNamespaces)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Debug("Getting deployment list done.")

	// log.Debug("Get statefulset list...")
	// statefulsets, statefulsetNodes, err := kubernetes.StatefulSetsOnNodes(blacklistedNamespaces)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Debug("Getting statefulset list done.")

	// log.Debug("Listing nodes...")
	// nodes, err := kubernetes.Nodes()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Debug("Listing nodes done.")

	// log.Debug("Map nodes to Kubernetes object")
	// clusterKubernetesObject, err := mapper.MapNodes(
	// 	viper.GetString("clustername"),
	// 	nodes,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Debug("Map deployments to Kubernetes objects")
	// deploymentKubernetesObjects := mapper.MapDeployments(viper.GetString(clusterNameFlag), deployments, deploymentNodes)
	// log.Debug("Map statefulsets to Kubernetes objects")
	// statefulsetKubernetesObjects := mapper.MapStatefulSets(viper.GetString(clusterNameFlag), statefulsets, statefulsetNodes)

	// kubernetesObjects := make([]mapper.KubernetesObject, 0)
	// kubernetesObjects = append(kubernetesObjects, *clusterKubernetesObject)
	// kubernetesObjects = append(kubernetesObjects, deploymentKubernetesObjects...)
	// kubernetesObjects = append(kubernetesObjects, statefulsetKubernetesObjects...)

	ldif := mapper.LDIF{
		ConnectorID:      viper.GetString(connectorIDFlag),
		ConnectorType:    "leanix-k8s-connector",
		ConnectorVersion: connectorVersion,
		LxVersion:        lxVersion,
		LxWorkspace:      viper.GetString(lxWorkspaceFlag),
		Description:      "Map Kubernetes objects to LeanIX Fact Sheets",
		Content:          kubernetesObjects,
	}

	log.Debug("Marshal ldif")
	ldifByte, err := storage.Marshal(ldif)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Upload %s to %s", storage.LdifFileName, viper.GetString("storage-backend"))
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
	uploader.Upload(ldifByte, debugLogBuffer.Bytes())
	log.Info("-----------End-----------")
}

func ServerPreferredListableResources(d discovery.DiscoveryInterface) ([]*metav1.APIResourceList, error) {
	all, err := discovery.ServerPreferredResources(d)
	return discovery.FilteredBy(discovery.ResourcePredicateFunc(func(groupVersion string, r *metav1.APIResource) bool {
		return strings.Contains(r.Verbs.String(), "list")
	}), all), err
}

func parseFlags() error {
	flag.String(clusterNameFlag, "", "unique name of the Kubernetes cluster")
	flag.String(storageBackendFlag, storage.FileStorage, fmt.Sprintf("storage where the %s file is placed (%s, %s)", storage.LdifFileName, storage.FileStorage, storage.AzureBlobStorage))
	flag.String(azureAccountNameFlag, "", "Azure storage account name")
	flag.String(azureAccountKeyFlag, "", "Azure storage account key")
	flag.String(azureContainerFlag, "", "Azure storage account container")
	flag.String(localFilePathFlag, ".", "path to place the ldif file when using local file storage backend")
	flag.Bool(verboseFlag, false, "verbose log output")
	flag.String(connectorIDFlag, "", "unique id of the LeanIX Kubernetes connector")
	flag.StringSlice(blacklistNamespacesFlag, []string{""}, "list of namespaces that are not scanned")
	flag.String(lxWorkspaceFlag, "", "name of the LeanIX workspace the data is sent to")
	flag.Bool(localFlag, false, "use local kubeconfig from home folder")
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
		return fmt.Errorf("%s flag must be set", clusterNameFlag)
	}
	if viper.GetString(connectorIDFlag) == "" {
		return fmt.Errorf("%s flag must be set", connectorIDFlag)
	}
	if viper.GetString(lxWorkspaceFlag) == "" {
		return fmt.Errorf("%s flag must be set", lxWorkspaceFlag)
	}
	return nil
}

// InitLogger initialise the logger for stdout and log file
func initLogger() (logging.LeveledBackend, *bytes.Buffer) {
	format := logging.MustStringFormatter(`%{time} ▶ [%{level:.4s}] %{message}`)
	logging.SetFormatter(format)

	// stdout logging backend
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutLeveled := logging.AddModuleLevel(stdout)

	// file logging backend
	var mem bytes.Buffer
	fileLogger := logging.NewLogBackend(&mem, "", 0)
	logging.SetBackend(fileLogger, stdoutLeveled)
	return stdoutLeveled, &mem
}

func enableVerbose(logger logging.LeveledBackend, verbose bool) {
	if verbose {
		logger.SetLevel(logging.DEBUG, "")
	} else {
		logger.SetLevel(logging.INFO, "")
	}
}
