node('go11') {
  stage('checkout') {
    sh """
      git rev-parse --is-inside-work-tree # timeout=10
      git config remote.origin.url https://github.com/tsungming/alameda.git # timeout=10
      git --version # timeout=10      
      git fetch --tags --progress https://github.com/tsungming/alameda.git +refs/heads/*:refs/remotes/origin/*
      git rev-parse origin/auto-pr1^{commit} # timeout=10
      ls -la
    """
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
