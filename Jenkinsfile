node('go11') {
  stage('checkout') {
        git branch: 'master', url: "https://github.com/containers-ai/alameda.git"
  }
  stage("Build Operator") {
    sh """
      pwd 
      ls -la ${env.WORKSPACE}
    """
  }
}
