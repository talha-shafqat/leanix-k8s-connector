pipeline {
    agent any
    node {
        // Install the desired Go version
        def root = tool name: 'Go 1.8', type: 'go'
 
        // Export environment variables pointing to the directory where Go was installed
        withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
            sh 'go version'
        }
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
                echo 'Here we need to run helm command to deploy to the leanix int cluster'
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