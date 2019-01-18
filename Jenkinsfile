node('go11') {
  stage('checkout') {
    sh """
    git clone https://github.com/tsungming/alameda.git
    """
  }
  stage('Example') {
    if (env.BRANCH_NAME == 'master') {
            echo 'I only execute on the master branch'
    } else {
            echo 'I execute elsewhere'
    }
  }
  stage("Build Operator") {
    sh """
      cd alameda
      git fetch --tags --progress https://github.com/tsungming/alameda.git +refs/heads/*:refs/remotes/origin/*
      git rev-parse origin/auto-pr1^{commit}
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
