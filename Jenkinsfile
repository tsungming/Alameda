node('go11') {
  stage('checkout') {
        git url: "https://github.com/tsungming/alameda.git", branch: 'master'
  }
  stage("Build Operator") {
    sh """
      cat Jenkinsfile
      pwd 
      ls -la ${env.WORKSPACE}
      echo "new branch"
    """
    pullRequest.addLabel('Build Failed')
    if (env.CHANGE_ID) {
      pullRequest.addLabel('Build Failed')
      echo "new pr1"
    }
  }
}
