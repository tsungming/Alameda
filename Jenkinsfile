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
      def getRepoURL() {
        sh "git config --get remote.origin.url > .git/remote-url"
        return readFile(".git/remote-url").trim()
      }
      
      def getCommitSha() {
        sh "git rev-parse HEAD > .git/current-commit"
        return readFile(".git/current-commit").trim()
      }
      
      def updateGithubCommitStatus(build) {
        // workaround https://issues.jenkins-ci.org/browse/JENKINS-38674
        repoUrl = getRepoURL()
        commitSha = getCommitSha()
      
        step([
          $class: 'GitHubCommitStatusSetter',
          reposSource: [$class: "ManuallyEnteredRepositorySource", url: repoUrl],
          commitShaSource: [$class: "ManuallyEnteredShaSource", sha: commitSha],
          errorHandlers: [[$class: 'ShallowAnyErrorHandler']],
          statusResultSource: [
            $class: 'ConditionalStatusResultSource',
            results: [
              [$class: 'BetterThanOrEqualBuildResult', result: 'SUCCESS', state: 'SUCCESS', message: build.description],
              [$class: 'BetterThanOrEqualBuildResult', result: 'FAILURE', state: 'FAILURE', message: build.description],
              [$class: 'AnyBuildResult', state: 'FAILURE', message: 'Loophole']
            ]
          ]
        ])
      }      
    }
  }
}