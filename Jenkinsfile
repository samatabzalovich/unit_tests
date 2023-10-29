pipeline {
  agent any
  stages {
    stage('Init') {
      steps {
        echo 'Initializing...'
        tool 'Go 1.8'
      }
    }

    stage('running') {
      steps {
        sh 'go run ./cmd/api'
      }
    }

  }
}