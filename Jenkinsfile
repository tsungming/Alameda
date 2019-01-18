node('go11') {
  stage('checkout') {
    sh """
      ls -la /var/lib/jenkins/jobs/tutorial-cicd/jobs/tutorial-cicd-alameda
    """"
    // git url: "https://github.com/tsungming/alameda.git", branch: 'auto-p1'
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
      echo "new pr2"
    }
  }
}
