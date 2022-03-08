#!/usr/bin/env groovy

@Library('jenkins_shared_library@main')_

pipeline {
    agent {
      label "golang-exec"
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
    }
    stages {
        stage('Pre Test') {
            steps {
                container('golang') {
                    echo 'Installing dependencies'
                    env.APP_DB_HOST = "10.46.142.201"
                    env.APP_DB_PORT = 3306
                    env.APP_DB_USERNAME = "user"
                    env.APP_DB_PASSWORD = "BcGH2Gj41J5VF1"
                    env.APP_DB_NAME = "ksar"
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
        stage('Build ksar') {
            steps {
                container('golang') {
                    echo 'Build Application'
                    sh 'go build -o ksar main.go'
                }
            }
        }
    }
}