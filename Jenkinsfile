node('go11') {
  stage('checkout') {
        git url: "https://github.com/tsungming/alameda.git", branch: 'auto1'
  }
  stage("Build Operator") {
    sh """
      pwd 
      ls -la ${env.WORKSPACE}      
      echo ${env.BRANCH_NAME}
      echo ${env.CHANGE_ID}
    """
    echo "Running ${env.BUILD_ID} on ${env.JENKINS_URL}"
    
    def printParams() {
      env.getEnvironment().each { name, value -> println "Name: $name -> Value $value" }
    }
    printParams()
  }
}
