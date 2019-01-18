node('go11') {
  stage('checkout') {
    sh """
    git clone https://github.com/tsungming/alameda.git
    """
  }
  stage("Build Operator") {
    sh """
      cd alameda
      cat Jenkinsfile
      pwd 
      ls -la ${env.WORKSPACE}
      echo "new branch"
    """
    pullRequest.addLabel('Build Failed')
    if (env.CHANGE_ID) {
      pullRequest.addLabel('Build Failed')
      echo "new pr2"
    }
  }
}
