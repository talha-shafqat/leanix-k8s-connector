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
                branch 'master'
            }
            steps {
                sh 'make'
            }
        }
    }
            }
        }
    }
}