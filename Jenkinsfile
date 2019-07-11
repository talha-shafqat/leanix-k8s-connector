pipeline {
    agent any

    stages {
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Build') {
            if(env.BRANCH_NAME == 'master'){
                steps {
                    sh 'make'
                }
            }
        }
    }
    post {
        success {
            script {
                if (env.CHANGE_ID) {
                    pullRequest.createStatus(status: 'success')
                }
            }
        }
        failure {
            script {
                if (env.CHANGE_ID) {
                    pullRequest.createStatus(status: 'failure')
                }
            }
        }
    }
}