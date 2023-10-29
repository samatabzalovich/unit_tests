pipeline {
  agent any
  stages {
    stage('Init') {
      parallel {
        stage('Init') {
          steps {
            echo 'Initializing...'
          }
        }

        stage('running') {
          steps {
            sh 'go run ./cmd/api'
          }
        }

      }
    }

  }
}