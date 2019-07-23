# LeanIX Kubernetes Connector
The LeanIX Kubernetes Connector collects information from Kubernetes.

## Overview
The LeanIX Kubernetes Connector runs in the Kubernetes cluster as container itself and collects information from the cluster, nodes, deployments, statefulsets and pods. Those informations are sanitized and brought into the LDIF (LeanIX Data Interchange Format) format that LeanIX understands. The output then is stored in the `kubernetes.ldif` file that gets imported into LeanIX.

## Getting started
Depending on how you would like to run the LeanIX Kubernetes Connector the installation steps differ. You can run the connector as a container in the Kubernetes cluster itself or executing the connector as CLI command directly on the Kubernetes nodes.

We recommend to deploy and run the connector as a container in the Kubernetes cluster.

### Architecture
The LeanIX Kubernetes Connector gets deployed via a Helm chart into the Kubernetes cluster as a CronJob. All necessary requirements like the ServiceAccount, the ClusterRole and ClusterRoleBinding are deployed also by the Helm chart.

Only necessary permissions are given to the connector as listed below and limited to get, list and watch operations.

|apiGroups   |resources                                         |verbs|
|------------|--------------------------------------------------|-----|
|""          |namespaces, nodes, pods, replicationcontrollers   |get, list, watch|
|"apps"      |daemonsets, deployments, replicasets, statefulsets|get, list, watch|
|"extensions"|daemonsets, deployments, replicasets              | get, list, watch|

The CronJob is configured to run every minute and spins up a new pod of the LeanIX Kubernetes Connector. As mentioned in the overview the connector creates the `kubernetes.ldif` file and logs into the `leanix-k8s-connector.log` file.

Two storage backend types are natively supported by the connector.

1. file
2. azureblob

The `file` storage backend lets you use every storage that can be provided to Kubernetes through a PersistentVolume and a PersistentVolumeClaim.

> When using the `file` storage backend you must pre-create the PersistentVolume and PersistentVolumeClaim the LeanIX Kubernetes Connector should use.

The `azureblob` storage backend leverages an Azure Storage account you must provide to store the `.ldif` and `.log` files

### Installation - Helm chart
Before you can install the LeanIX Kubernetes Connector make sure that the following pre-requisites are fulfilled on your local workstation.

- [kubectl is installed](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [helm client is installed](https://helm.sh/docs/using_helm/#installing-the-helm-client)
- [git is installed](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

Run `git clone https://github.com/leanix/leanix-k8s-connector.git` to clone the leanix-k8s-connector repository to your local workstation.

#### `file` storage backend
tbd

#### `azureblob` storage backend
tbd

### Installation - CLI
tbd

## Known issues
tbd

## Version history

|Connector version  |Integration version  |Helm chart version  |
|:-----------------:|:-------------------:|:------------------:|
|0.0.1              |3                    |0.0.1

## Roadmap
- [ ] Collect information from DaemonSets