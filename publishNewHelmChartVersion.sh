#!/bin/bash

set -e

cd ./helm

# Creates or updates a new Helm chart package
helm package ./leanix-k8s-connector

# Creates or updates the Helm chart repository index
helm repo index .