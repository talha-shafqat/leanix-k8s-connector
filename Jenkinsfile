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
}