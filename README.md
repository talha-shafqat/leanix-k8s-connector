# LeanIX Kubernetes Connector
The LeanIX Kubernetes Connector collects information from Kubernetes.

## Overview
The LeanIX Kubernetes Connector runs in the Kubernetes cluster as a container itself and collects information from the cluster, nodes, deployments, statefulsets and pods. Those informations are sanitized and brought into the LDIF (LeanIX Data Interchange Format) format that LeanIX understands. The output then is stored in the `kubernetes.ldif` file that gets imported into LeanIX.

## Getting started
Depending on how you would like to run the LeanIX Kubernetes Connector the installation steps differ. You can run the connector as a container in the Kubernetes cluster itself or executing the connector as CLI command directly on the Kubernetes nodes.

We recommend deploying and running the connector as a container in the Kubernetes cluster.

### Architecture
The LeanIX Kubernetes Connector gets deployed via a Helm chart into the Kubernetes cluster as a CronJob. All necessary requirements like the ServiceAccount, the ClusterRole and ClusterRoleBinding are deployed also by the Helm chart.

Only necessary permissions are given to the connector as listed below and limited to get, list and watch operations.

|apiGroups   |resources                                         |verbs           |
|------------|--------------------------------------------------|----------------|
|""          |namespaces, nodes, pods, replicationcontrollers   |get, list, watch|
|"apps"      |daemonsets, deployments, replicasets, statefulsets|get, list, watch|
|"extensions"|daemonsets, deployments, replicasets              |get, list, watch|

The CronJob is configured to run every minute and spins up a new pod of the LeanIX Kubernetes Connector. As mentioned in the overview the connector creates the `kubernetes.ldif` file and logs into the `leanix-k8s-connector.log` file.

Currently two storage backend types are natively supported by the connector.

- file
- azureblob

The `file` storage backend lets you use every storage that can be provided to Kubernetes through a PersistentVolume and a PersistentVolumeClaim.

> When using the `file` storage backend you must pre-create the PersistentVolume and PersistentVolumeClaim the LeanIX Kubernetes Connector should use.

The `azureblob` storage backend leverages an Azure Storage account you must provide to store the `.ldif` and `.log` files.

### Installation - Helm chart
Before you can install the LeanIX Kubernetes Connector make sure that the following pre-requisites are fulfilled on your local workstation.

- [kubectl is installed](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [helm client is installed](https://helm.sh/docs/using_helm/#installing-the-helm-client)
- [git is installed](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

Run `git clone https://github.com/leanix/leanix-k8s-connector.git` to clone the leanix-k8s-connector repository to your local workstation.

#### `file` storage backend
The first step to get started with the `file` storage backend type is to create the PersistentVolume and PersistentVolumeClaim in advance.

In the following example the creation of a PV and PVC to connect to Azure Files is shown.

Start with the creation of an Azure Storage account and an Azure file share as described in the Azure documentation. In our example we used _leanixk8sconnector_ as file share name.

1. [Create a storage account](https://docs.microsoft.com/en-us/azure/storage/common/storage-quickstart-create-account?tabs=azure-portal)
2. [Create a file share in Azure Files](https://docs.microsoft.com/en-us/azure/storage/files/storage-how-to-create-file-share)

Next, create a Kubernetes secret with the Azure Storage account name and the Azure Storage account key. The information about the name and the key can be retrieved directly via the Azure portal.

- [Access keys](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-manage#access-keys)

```Bash
kubectl create secret generic azure-secret --from-literal=azurestorageaccountname={STORAGE_ACCOUNT_NAME} --from-literal=azurestorageaccountkey={STORAGE_KEY}
```

Afterwards create the PV and PVC using the template below running the `kubectl apply -f template.yaml` command.

```YAML
apiVersion: v1
kind: PersistentVolume
metadata:
  name: azurefile
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  azureFile:
    secretName: azure-secret
    shareName: leanixk8sconnector
    readOnly: false
  mountOptions:
  - dir_mode=0777
  - file_mode=0777
  - uid=1000
  - gid=1000
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: azurefile
spec:
  accessModes:
    - ReadWriteMany
  volumeName: azurefile
  storageClassName: ""
  resources:
    requests:
      storage: 1Gi
```

Run `kubectl get pv && kubectl get pvc` and check your output. It should look like this.

```
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
azurefile   1Gi        RWX            Retain           Bound    default/azurefile                           45s
NAME        STATUS   VOLUME      CAPACITY   ACCESS MODES   STORAGECLASS   AGE
azurefile   Bound    azurefile   1Gi        RWX                           45s
```

Finally we use the Helm chart deploying the LeanIX Kubernetes Connector to the Kubernetes cluster.

The following command deploys the connector to the Kubernetes cluster and overwrites the parameters in the `values.yaml` file.

|Parameter          |Default value|Provided value      |Notes                                                                                        |
|-------------------|-------------|--------------------|---------------------------------------------------------------------------------------------|
|clustername        |kubernetes   |aks-cluster         |The name of the Kubernetes cluster.                                                          |
|claimName          |""           |azurefile           |The name of the PVC used to store the `kubernetes.ldif` and `leanix-k8s-connector.log` files.|
|connectorID        |Random UUID  |aks-cluster         |The name of the Kubernetes cluster. If not provided a random UUID is generated per default.  |
|blacklistNameSpaces|kube-system  |kube-system, default|Namespaces that are not scanned by the connector. Must be provided in the format `"{kube-system,default}"` when using the `--set` option|
|verbose            |false        |true                |Enables verbose logging on the stdout interface of the container                             |

```Bash
helm upgrade --install leanix-k8s-connector ./helm/leanix-k8s-connector \
--set args.clustername=aks-cluster \
--set args.file.claimName=azurefile \
--set args.connectorID=aks-cluster \
--set args.blacklistNamespaces="{kube-system,default}" \
--set args.verbose=true
```

Beside the option to override the default values and provide values via the `--set` option of the `helm` command, you can also edit the `values.yaml` file.

```YAML
...

args:
  clustername: aks-cluster
  connectorID: aks-cluster
  verbose: true
  storageBackend: file
  file:
    localFilePath: "/mnt/leanix-k8s-connector"
    claimName: "azurefile"
  azureblob:
    secretName: ""
    container: ""
  blacklistNamespaces:
  - "kube-system"
  - "default"

...
```

#### `azureblob` storage backend
tbd

### Installation - CLI
tbd

## Known issues
tbd

## Version history

|Connector version  |Integration version  |Helm chart version  |
|:-----------------:|:-------------------:|:------------------:|
|0.0.1              |3                    |0.0.1               |

## Roadmap
- [ ] Collect information from DaemonSets