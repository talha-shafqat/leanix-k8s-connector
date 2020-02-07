# LeanIX Kubernetes Connector Changelog

## Release 2020-02-07 - 2.0.0-beta2

### Release Notes

* New Features
  * Enable information extraction for the following Kubernetes resources:
    * replicasets
    * replicationcontrollers

* Changes
  * The `connectorId` field gets pinned to `Kubernetes` in the generated LDIF file.
  * The customer provided value for the `connectorId` field is moved to the section `customFields` and is mapped to the field `connectorInstance` in the generated LDIF file.

## Release 2020-01-14 - 2.0.0-beta1

### Release Notes

* New Features
  * Enable information extraction for the following Kubernetes resources:
    * serviceaccounts
    * services
    * nodes
    * pods
    * namespaces
    * configmaps
    * persistentvolumes
    * persistentvolumeclaims
    * deployments
    * statefulsets
    * daemonsets
    * customresourcedefinitions
    * clusterrolebindings
    * rolebindings
    * clusterroles
    * roles
    * ingresses
    * networkpolicies
    * horizontalpodautoscalers
    * podsecuritypolicies
    * storageclasses
    * cronjobs
    * jobs

* Removed Features
  * Custom Kubernetes resource extraction and information aggregation for:
    * deployments
    * statefulsets

* Changes
  * Permission change for the `leanix-k8s-connector` service account

## Release 2019-09-26 - 1.1.0

### Release Notes

* New Features
  * Enable wildcards for namespace filtering
