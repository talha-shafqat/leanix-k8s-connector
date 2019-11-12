pipeline {
    agent any
    tools {
        go '1.12.7'
    }
    environment {
        VERSION = """${sh(
                returnStdout: true,
                script: 'make version'
            )}"""
    }
    stages {
        stage('Test') {
            steps {
                sh 'make'
                sh 'make image'
                sh 'docker run leanix/leanix-k8s-connector:${VERSION} --help | grep "pflag: help requested" '
                sh 'docker push leanix/leanix-k8s-connector:${VERSION}'
                sh 'helm upgrade --install leanix-k8s-connector ./helm/leanix-k8s-connector --set image.tag=${VERSION} --set args.clustername=leanix-cluster --set args.storageBackend=azureblob --set args.azureblob.secretName=azure-secret --set args.azureblob.container=connector --set args.connectorID=leanix-int --set args.lxWorkspace=leanix --set args.verbose=true'
            }
        }
        stage('Build') {
            when {
                anyOf {
                    branch 'master'
                    branch 'develop'
                }
            }
            steps {
                sh 'make'
                sh 'make image'
                sh 'make push'
            }
        }
        stage('Deploy to int cluster') {
            when {
                anyOf {
                    branch 'master'
                    branch 'develop'
                }
            }
            steps {
                sh 'helm upgrade --install leanix-k8s-connector ./helm/leanix-k8s-connector --set image.tag=${VERSION} --set args.clustername=leanix-cluster --set args.storageBackend=azureblob --set args.azureblob.secretName=azure-secret --set args.azureblob.container=connector --set args.connectorID=leanix-int --set args.lxWorkspace=leanix'
            }
        }
        // stage('Release approval'){
        //     when {
        //         branch 'master'
        //     }
        //     input "Release new version?"
        // }
        // stage('Release') {
        //     when {
        //         branch 'master'
        //     }
        //     steps {
        //         echo 'Set the version variable as default for image tag in helm chart'
        //     }
        // }
    }
}