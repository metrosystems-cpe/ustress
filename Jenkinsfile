pipeline {
  environment {
    SLACK_CHANNEL = '@ionut.ilie,@bogdanb'
  }
  agent {
    kubernetes {
      label 'dockerdind'
      yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
  vertical: "reliability"
  ci: true
spec:
  containers:
  - name: docker
    image: docker:dind
    command: ["sh"]
    args: ['-c', 'apk add -U bash && dockerd --host=unix:///var/run/docker.sock --host=tcp://0.0.0.0:2375 --mtu=800']
    env:
    - name: DOCKER_REGISTRY_HOST
      valueFrom:
        configMapKeyRef:
          name: docker-1
          key: host
    - name: DOCKER_USERNAME
      valueFrom:
        secretKeyRef:
          name: docker-1
          key: username
    - name: DOCKER_PASSWORD
      valueFrom:
        secretKeyRef:
          name: docker-1
          key: password
    - name: IDAM_CLIENT_ID
      valueFrom:
        secretKeyRef:
          name: idam
          key: client_id
    - name: IDAM_SECRET
      valueFrom:
        secretKeyRef:
          name: idam
          key: password
    - name: SLACK_TOKEN
      valueFrom:
        secretKeyRef:
          name: slack
          key: token

    tty: true
    securityContext:
      privileged: true
"""
    }
  }
  stages {
    stage('build') {
      steps {
        container('docker') {
          sh '''#!/bin/bash
          # source ./ci/utils.sh
          # alpine_init
          # docker build -t build-image - < ./ci/build.Dockerfile
          # docker run --rm -t -v $PWD:/mnt/workspace build-image
          # docker build --network host -t restmonkey:latest -f ./ci/build.Dockerfile .
          apk add -U curl git ca-certificates make
          make docker
          make docker-release
          '''
        }
      }
    }
    stage('deploy pp') {
      steps{
        container('docker') {
          sh '''#!/bin/bash
          echo "
          deployment is not configured for this app
          delete the old pod from kubernetes and the app will update
          "  

          # source ./ci/utils.sh
          # remote_image_url=$(build_and_tag_image "pp" "reliability" "restmonkey" "./ci/run-restmonkey.Dockerfile")
          # remote_image_url=$(tag_image "pp" "reliability" "restmonkey")
      
          # printf "%-7s: %s %s \n" "INFO" "docker target image:tag" ${remote_image_url}
          # [[ "$remote_image_url" = "" ]] && exit 1
          # push_image_to_registry "$remote_image_url"
          # deploy_to_ds "$remote_image_url" "pp" "./ci/restmonkey-ds-payload.json" "./ci/restmonkey-resources.json"
          '''
        }
      }
    }
  }
  post {
    failure {
      container('docker') {
        script {
          def token = sh (
              script: 'echo -n ${SLACK_TOKEN}',
              returnStdout: true
          ).trim()
          def teamDomain = "metro-dr"
          def msg = "FAILURE: ${env.JOB_NAME} #${env.BUILD_NUMBER}"
          slackSend(color: "#FF9FA1", message: msg, channel: "${env.SLACK_CHANNEL}", teamDomain: teamDomain, token: token, tokenCredentialId: "slack-token")
          mail to: '',
            subject: "Failed pipeline: ${currentBuild.fullDisplayName}",
            body: "Something is wrong with ${env.BUILD_URL}"
        }
      }
    }
  }
}
