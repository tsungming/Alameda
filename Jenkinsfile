node('go11') {
  stage('checkout') {
        git url: "https://github.com/tsungming/alameda.git", branch: 'master'
  }
  stage("Build Operator") {
    sh """
      pwd 
      ls -la ${env.WORKSPACE}
    """
  }
}
