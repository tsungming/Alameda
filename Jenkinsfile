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
                sh '''
                  go env
                  mkdir -p /go/src/github.com/containers-ai/alameda
                  cp -R . /go/src/github.com/containers-ai/alameda
                '''
              }
            }
          }
        }
      }
    }
    stage("Build Operator") {
      steps {        
        sh '''
          cd /go/src/github.com/containers-ai/alameda/operator
          make manager
        '''
      }
      post {
        always {
          step([
            $class: 'GitHubCommitStatusSetter',
            contextSource: [$class: 'ManuallyEnteredCommitContextSource', context: 'jenkins-build-operator']
          ])
        }
        failure {
          step([
            $class: 'GitHubCommitStatusSetter',
            contextSource: [$class: 'ManuallyEnteredCommitContextSource', context: 'jenkins-build-operator'],
            statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: "message", state: "FAILURE"]]]
          ])
        }
        success {
          step([
            $class: 'GitHubCommitStatusSetter',
            contextSource: [$class: 'ManuallyEnteredCommitContextSource', context: 'jenkins-build-operator'],
            statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: "message", state: "SUCCESS"]]]
          ])
        }
      }
    }
    stage("Build Operator") {
      steps {
        sh '''
          cd /go/src/github.com/containers-ai/alameda/datahub
          make datahub
        '''
      }
    }
    stage("Test Operator") {
      steps {        
        sh '''
          cd /go/src/github.com/containers-ai/alameda/operator
          make test
        '''
      }
    }
    stage("Test Datahub") {
      steps {        
        sh '''
          cd /go/src/github.com/containers-ai/alameda/datahub
          make test          
        '''
      }
    }
  }
  post {
    always {
      step([
        $class: 'GitHubCommitStatusSetter',
        contextSource: [$class: 'ManuallyEnteredCommitContextSource', context: 'jenkins-ci']
      ])
    }
    failure {
      step([
        $class: 'GitHubCommitStatusSetter',
        contextSource: [$class: 'ManuallyEnteredCommitContextSource', context: 'jenkins-ci'],
        statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: "message", state: "FAILURE"]]]
      ])
    }
    success {
      step([
        $class: 'GitHubCommitStatusSetter',
        contextSource: [$class: 'ManuallyEnteredCommitContextSource', context: 'jenkins-ci'],
        statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: "message", state: "SUCCESS"]]]
      ])
    }
  }
}