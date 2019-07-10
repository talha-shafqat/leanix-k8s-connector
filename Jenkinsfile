pipeline {
    agent any

    stages {
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Build') {
            steps {
                echo 'Building...'
            }
        }
    }
    post {
        success {
            githubNotify status: 'SUCCESS'
        }
        failure {
            githubNotify status: 'FAILURE'
        }
    }
}