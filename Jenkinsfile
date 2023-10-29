pipeline {
  agent {
    docker {
      image 'golang'
    }

  }
  stages {
    stage('running') {
      steps {
        sh 'go run ./cmd/api'
      }
    }

  }
}