#!/usr/bin/env groovy

@Library('jenkins_shared_library@main')_

pipeline {
  agent {
    kubernetes {
      yaml '''
        apiVersion: v1
        kind: Pod
        metadata:
          labels:
            jenkins-label: golang-docker
        spec:
          containers:
          - name: golang
            image: golang:1.17.8
            command:
            - sleep
            args:
            - 99d
            tty: true
          - name: kaniko
            image: gcr.io/kaniko-project/executor:debug
            command:
            - sleep
            args:
            - 9999999
            tty: true
            volumeMounts:
            - name: kaniko-secret
              mountPath: /kaniko/.docker
          restartPolicy: Never
          volumes:
          - name: kaniko-secret
            secret:
                secretName: dockercred
                items:
                - key: .dockerconfigjson
                  path: config.json
        '''
    }
  }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
    }
  stages {
        stage('Checkout Source') {
            steps {
                git url:'https://github.com/kumarsarath588/ksar.git', branch:'main'
            }
        }
        stage('Set Test params') {
            steps {
                script {
                    env.APP_DB_HOST = "10.46.142.201"
                    env.APP_DB_PORT = 3306
                    env.APP_DB_USERNAME = "user"
                    env.APP_DB_PASSWORD = "BcGH2Gj41J5VF1"
                    env.APP_DB_NAME = "ksar"
                }
            }
        }
        stage('Install go dependencies') {
            steps {
                container('golang') {
                    echo 'Installing dependencies'
                    sh 'go version'
                    sh 'go mod tidy'
                    sh 'go mod vendor'
                }
            }
        }
        stage('Testing') {
            steps {
                container('golang') {
                    echo 'Testing Application'
                    sh 'go test -v'
                }
            }
        }
        stage('Build ksar Application') {
            steps {
                container('golang') {
                    echo 'Build Application'
                    sh 'go build -o ksar main.go'
                }
            }
        }
        stage('Build ksar Docker image') {
            steps {
                container('kaniko') {
                    echo 'Build Docker Image'
                    sh '/kaniko/executor --context `pwd` --destination kumarsarath588/ksar:1.0.0'
                }
            }
        }
    }
}