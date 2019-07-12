pipeline {
    agent any
    tools { 
        go '1.12.7' 
    }
    environment {
        AZURE_STORAGE_ACCOUNT_NAME = 'mastest534'
        AZURE_STORAGE_ACCOUNT_KEY = credentials('mas-azure-storage-account-key')
    }
    stages {
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Build') {
            when { 
                anyOf { 
                    branch 'master'
                    branch 'develop'
                    branch 'feature/HOOK-426-jenkins-job'
                } 
            }
            steps {
                sh 'make'
                sh 'make image'
                sh 'make push'
                v = sh (script: 'make version', returnStdout: true).trim()
                env.version = v
            }
        }
        stage('Deploy to int cluster') {
            when { 
                anyOf { 
                    branch 'master'
                    branch 'develop'
                    branch 'feature/HOOK-426-jenkins-job'
                } 
            }
            steps {
                sh 'helm upgrade --install leanix-k8s-connector ./helm/leanix-k8s-connector --set image.tag=${env.version} --set args.clustername=leanix-westeurope-int --set args.storageBackend=azureblob --set args.azureblob.accountKey=$(AZURE_STORAGE_ACCOUNT_KEY) --set args.azureblob.accountName=$(AZURE_STORAGE_ACCOUNT_NAME) --set args.azureblob.container=connector --set args.connectorID=leanix-int'
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