node('go11') {
  echo sh(returnStdout: true, script: 'env')
  stage('checkout') {
        git url: "https://github.com/tsungming/alameda.git", branch: 'auto1'
  }
  stage("Build Operator") {
    sh """
      pwd 
      ls -la ${env.WORKSPACE}      
      echo ${env.BRANCH_NAME}
      echo ${env.CHANGE_ID}
      env
    """
  }
}
