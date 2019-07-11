pipeline {
    agent any

    stages {
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Build') {
            when {
                branch 'develop'
            }
            steps {
                sh 'make'
                sh 'make image'
                sh 'make push'
            }
        }
    }
}