pipeline {
    agent any
    stages {
        stage('Checkout') {
            steps {
                git branch: '${BRANCH_NAME}', url: 'YOUR_GIT_REPOSITORY_URL'
            }
        }
        stage('Build Backend') {
            steps {
                script {
                    if (BRANCH_NAME != 'frontend') {
                        sh 'docker-compose build backend'
                    }
                }
            }
        }
        stage('Build Frontend') {
             steps {
                script {
                    if (BRANCH_NAME != 'backend') {
                        sh 'docker-compose build frontend'
                    }
                }
            }
        }
        stage('Run Tests') {
            steps {
                script {
                    if (BRANCH_NAME != 'frontend') {
                        sh 'docker-compose run backend go test ./...'
                    }
                    if (BRANCH_NAME != 'backend') {
                        sh 'docker-compose run frontend npm run test:unit' // or test:e2e
                    }
                }
            }
        }
        stage('Database Migrations') {
            steps {
                script {
                    if (BRANCH_NAME != 'frontend' && BRANCH_NAME == 'develop') { // Only on develop branch
                        sh 'docker-compose run backend go run ./migrations/migrate.go up'
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    if (BRANCH_NAME == 'main') {
                        // Deploy to production environment (e.g., Kubernetes, AWS ECS)
                        sh 'docker-compose up -d' // Replace with your deployment commands
                    } else if (BRANCH_NAME == 'develop') {
                        // Deploy to staging environment
                        sh 'docker-compose up -d' // Replace with your deployment commands
                    }
                }
            }
        }
    }
    environment {
        BRANCH_NAME = "${env.BRANCH_NAME}"
    }
}