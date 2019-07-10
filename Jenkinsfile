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
            script {
                pullRequest.createStatus(status: 'success')
            }
        }
        failure {
            script {
                pullRequest.createStatus(status: 'failure')
            }
        }
    }
}