# LeanIX Kubernetes Connector

The LeanIX Kubernetes Connector collects information from Kubernetes.

## Table of contents

- [Overview](#Overview)
- [Getting started](#Getting-started)
  - [Architecture](#Architecture)
  - [Installation Helm chart](#Installation---Helm-chart)
    - [file storage backend](#file-storage-backend)
    - [azureblob storage backend](#azureblob-storage-backend)
- [Known issues](#Known-issues)
- [Version history](#Version-history)
- [Roadmap](#Roadmap)

## Overview

The LeanIX Kubernetes Connector runs in the Kubernetes cluster as a container itself and collects information from the cluster, nodes, deployments, statefulsets and pods. Those informations are sanitized and brought into the LDIF (LeanIX Data Interchange Format) format that LeanIX understands. The output then is stored in the `kubernetes.ldif` file that gets imported into LeanIX.

## Getting started

Depending on how you would like to run the LeanIX Kubernetes Connector the installation steps differ a bit and depends on the selected storage backend.

### Architecture

The LeanIX Kubernetes Connector gets deployed via a Helm chart into the Kubernetes cluster as a CronJob. All necessary requirements like the ServiceAccount, the ClusterRole and ClusterRoleBinding are deployed also by the Helm chart.

Only necessary permissions are given to the connector as listed below and limited to get, list and watch operations.

|apiGroups   |resources                                         |verbs           |
|------------|--------------------------------------------------|----------------|
|""          |namespaces, nodes, pods, replicationcontrollers   |get, list, watch|
|"apps"      |daemonsets, deployments, replicasets, statefulsets|get, list, watch|
|"extensions"|daemonsets, deployments, replicasets              |get, list, watch|

The CronJob is configured to run every minute and spins up a new pod of the LeanIX Kubernetes Connector. As mentioned in the overview the connector creates the `kubernetes.ldif` file and logs into the `leanix-k8s-connector.log` file.

Currently, two storage backend types are natively supported by the connector.

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

On the server-side the Helm server component Tiller must be deployed into the Kubernetes cluster.

- [Installing Tiller](https://helm.sh/docs/using_helm/#installing-tiller)

Run `git clone https://github.com/leanix/leanix-k8s-connector.git` to clone the leanix-k8s-connector repository to your local workstation and change the directory.

#### file storage backend

The first step to get started with the `file` storage backend type is to create the PersistentVolume and PersistentVolumeClaim in advance.

In the following example the creation of a PV and PVC to connect to Azure Files is shown.

Start with the creation of an Azure Storage account and an Azure file share as described in the Azure documentation. In our example we used _leanixk8sconnector_ as file share name.

1. [Create a storage account](https://docs.microsoft.com/en-us/azure/storage/common/storage-quickstart-create-account?tabs=azure-portal)
2. [Create a file share in Azure Files](https://docs.microsoft.com/en-us/azure/storage/files/storage-how-to-create-file-share)

Next, create a Kubernetes secret with the Azure Storage account name and the Azure Storage account key. The information about the name and the key can be retrieved directly via the Azure portal.

- [Access keys](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-manage#access-keys)

``` bash
kubectl create secret generic azure-secret --from-literal=azurestorageaccountname={STORAGE_ACCOUNT_NAME} --from-literal=azurestorageaccountkey={STORAGE_KEY}
```

Afterwards create the PV and PVC using the template below running the `kubectl apply -f template.yaml` command.

``` yaml
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

``` bash
NAME        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
azurefile   1Gi        RWX            Retain           Bound    default/azurefile                           45s
NAME        STATUS   VOLUME      CAPACITY   ACCESS MODES   STORAGECLASS   AGE
azurefile   Bound    azurefile   1Gi        RWX                           45s
```

Finally we use the Helm chart deploying the LeanIX Kubernetes Connector to the Kubernetes cluster.

The following command deploys the connector to the Kubernetes cluster and overwrites the parameters in the `values.yaml` file.

|Parameter          |Default value            |Provided value      |Notes                                                                                        |
|-------------------|-------------------------|--------------------|---------------------------------------------------------------------------------------------|
|clustername        |kubernetes               |aks-cluster         |The name of the Kubernetes cluster.                                                          |
|connectorID        |Random UUID              |aks-cluster         |The name of the Kubernetes cluster. If not provided a random UUID is generated per default.  |
|verbose            |false                    |true                |Enables verbose logging on the stdout interface of the container.
|storageBackend     |file                     |                    |The default value for the storage backend is `file`, if not provided.                        |
|localFilePath      |/mnt/leanix-k8s-connector|                    |The path that is used for mounting the PVC into the container and storing the `kubernetes.ldif` and `leanix-k8s-connector.log` files.|
|claimName          |""                       |azurefile           |The name of the PVC used to store the `kubernetes.ldif` and `leanix-k8s-connector.log` files.|
|blacklistNameSpaces|kube-system              |kube-system, default|Namespaces that are not scanned by the connector. Must be provided in the format `"{kube-system,default}"` when using the `--set` option|

``` bash
helm upgrade --install leanix-k8s-connector ./helm/leanix-k8s-connector \
--set args.clustername=aks-cluster \
--set args.connectorID=aks-cluster \
--set args.verbose=true \
--set args.file.claimName=azurefile \
--set args.blacklistNamespaces="{kube-system,default}"
```

Beside the option to override the default values and provide values via the `--set` option of the `helm` command, you can also edit the `values.yaml` file.

``` yaml
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

#### azureblob storage backend

The first step to get started with the `azureblob` storage backend type is to create an Azure Storage account and a container as described in the Azure documentation. In our example we used _leanixk8sconnector_ as container name.

1. [Create a storage account](https://docs.microsoft.com/en-us/azure/storage/common/storage-quickstart-create-account?tabs=azure-portal)
2. [Create a container](https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-portal#create-a-container)

Next, create a Kubernetes secret with the Azure Storage account name and the Azure Storage account key. The information about the name and the key can be retrieved directly via the Azure portal.

- [Access keys](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-manage#access-keys)

``` bash
kubectl create secret generic azure-secret --from-literal=azurestorageaccountname={STORAGE_ACCOUNT_NAME} --from-literal=azurestorageaccountkey={STORAGE_KEY}
```

Finally we use the Helm chart deploying the LeanIX Kubernetes Connector to the Kubernetes cluster.

The following command deploys the connector to the Kubernetes cluster and overwrites the parameters in the `values.yaml` file.

|Parameter          |Default value            |Provided value      |Notes                                                                                              |
|-------------------|-------------------------|--------------------|---------------------------------------------------------------------------------------------------|
|clustername        |kubernetes               |aks-cluster         |The name of the Kubernetes cluster.                                                                |
|connectorID        |Random UUID              |aks-cluster         |The name of the Kubernetes cluster. If not provided a random UUID is generated per default.        |
|verbose            |false                    |true                |Enables verbose logging on the stdout interface of the container                                   |
|storageBackend     |file                     |azureblob           |The default value for the storage backend is `file`, if not provided.                              |
|secretName         |""                       |azure-secret        |The name of the Kubernetes secret containing the Azure Storage account credentials.                |
|container          |""                       |leanixk8sconnector  |The name of the container used to store the `kubernetes.ldif` and `leanix-k8s-connector.log` files.|
|blacklistNameSpaces|kube-system              |kube-system, default|Namespaces that are not scanned by the connector. Must be provided in the format `"{kube-system,default}"` when using the `--set` option|

``` bash
helm upgrade --install leanix-k8s-connector ./helm/leanix-k8s-connector \
--set args.clustername=aks-cluster \
--set args.connectorID=aks-cluster \
--set args.verbose=true \
--set args.storageBackend=azureblob \
--set args.azureblob.secretName=azure-secret \
--set args.azureblob.container=leanixk8sconnector
--set args.blacklistNamespaces="{kube-system,default}"
```

Beside the option to override the default values and provide values via the `--set` option of the `helm` command, you can also edit the `values.yaml` file.

``` yaml
...

args:
  clustername: aks-cluster
  connectorID: aks-cluster
  verbose: true
  storageBackend: azureblob
  file:
    localFilePath: "/mnt/leanix-k8s-connector"
    claimName: ""
  azureblob:
    secretName: "azure-secret"
    container: "leanixk8sconnector"
  blacklistNamespaces:
  - "kube-system"
  - "default"

...
```

## Known issues

When the LeanIX Kubernetes Connector pod resides in an `Error` or `CrashLoopBackOff` state and you issued a `helm upgrade --install` command to fix it, you still the see the same pod instead of a new one.

This is not an issue of the LeanIX Kubernetes Connector itself. Instead it takes several minutes in this case until the `CronJob` creates a new pod.

If you do not want to wait until Kubernetes fix it itself, you can just delete the `Job` object.

Run `kubectl get jobs.batch` and look for the `Job` object with COMPLETIONS 0/1.

``` bash
NAME                              COMPLETIONS   DURATION   AGE
leanix-k8s-connector-1563961200   0/1           20m        20m
```

Issue `kubectl delete jobs.batch leanix-k8s-connector-1563961200` and you should see a new pod coming up afterwards.

## Version history

|Connector version  |Integration version  |Helm chart version  |
|:-----------------:|:-------------------:|:------------------:|
|0.0.1              |3                    |0.0.1               |
