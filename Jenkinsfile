pipeline {
  agent {
    node {
      // spin up a node.js slave pod to run this build on
      label 'go11'
    }
  }
  options {
    // set a timeout of 20 minutes for this pipeline
    timeout(time: 20, unit: 'MINUTES')
  }
  stages {
    stage('preamble') {
      steps {
        script {
          openshift.withCluster() {
            openshift.withProject() {
              echo "Using project: ${openshift.project()}"
              echo "Using project: ${env.GIT_COMMIT}"
              echo "Using project: ${env.GIT_BRANCH}"
            }
          }
        }
      }
    }
    stage('checkout') {
      steps {
        script {
          openshift.withCluster() {
            openshift.withCredentials() {
              openshift.withProject() {
                checkout([$class           : 'GitSCM',
                          branches         : [[name: "*/*"]],
                          userRemoteConfigs: [[url: "https://github.com/tsungming/alameda.git"]]
                ]);
                dir("${WORKSPACE}") {
                  def commit_id = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
                        echo "${commit_id}"
                }
              }
            }
          }
        }
      }
    }
    stage('build') {
      steps {
        echo "Perform Build"
        statusVerifier: allowRunOnStatus('SUCCESS')
      }
      post {
        always {
          githubPRStatusPublisher buildMessage: message(failureMsg: githubPRMessage('Can\'t set status; build failed.'), successMsg: githubPRMessage('Can\'t set status; build succeeded.')), errorHandler: statusOnPublisherError('UNSTABLE'), statusMsg: githubPRMessage('${GITHUB_PR_COND_REF} run ended'), statusVerifier: allowRunOnStatus('SUCCESS'), unstableAs: 'FAILURE'
                }
      }
    }
  }
}