pipeline {

    environment{
      REGISTRY = "raniamidaoui"
      IMAGE_NAME = "gomart-order-service"
      IMAGE_VERSION = "latest"
      //IMAGE_TEST= "test-order-test"
      LOGIN=credentials("sonar_login")
      NETWORK="tools_jenkins"
      SONARQUBE_URL="http://sonarqube:9000"
      PROJECT_KEY="gomart-microservice"
      K8S_CONFIG_NAME="kubeconfig"
      REGISTRY_CRED = 'docker_token'
      K8S_NAMESPACE = 'default'
      K8S_DEPLOYMENT_NAME = "order-microservice"
    }

    tools {
      go 'go_120'
    }

    agent any
    
    stages {
        
        stage('SonarTests') {
            steps{
                sh "docker network ls"
                script{
                    docker.image('sonarsource/sonar-scanner-cli').inside('-v /var/run/docker.sock:/var/run/docker.sock --entrypoint="" --net ${NETWORK}') {
                        sh "/opt/sonar-scanner/bin/sonar-scanner -Dsonar.projectKey=${PROJECT_KEY} -Dsonar.sources=. -Dsonar.host.url=${SONARQUBE_URL} -Dsonar.login=${LOGIN}"
                    }
                }
            }
        }


        // stage('Test') {
        //     steps {
        //         echo "Unit testing"
        //         sh "docker build --target test -t ${IMAGE_TEST} ."
        //     }
        //     post{
        //         success{
        //             echo "Unit Testing Passed"
        //         }
        //         failure{
        //             echo "Unit Testing Failed"
        //         }
        //     }
        // }

        stage("Build"){
            steps {
                echo "Docker Build"
                script{
                    dockerImg = docker.build("${REGISTRY}/${IMAGE_NAME}:${IMAGE_VERSION}", ".") 
                }
            }
            post{
                success{
                    echo "Build Successful"
                }
                failure{
                    echo "Build Failed"
                }
            }
        }
         
        stage("Push"){
            steps{
                script{
                    docker.withRegistry('', "${REGISTRY_CRED}") {
                        dockerImg.push()
                    }
                }
            }
            post{
                    success{
                        echo "Push Successful"
                    }
                    failure{
                        echo "Push Failed"
                    }
            }
        }
        
        stage('Restart Helm Chart') {
            steps {
                script{
                    docker.image('bitnami/kubectl').inside("--entrypoint=''  --net minikube"){
                        withKubeConfig([credentialsId: "$K8S_CONFIG_NAME", namespace: "$K8S_NAMESPACE"]) {
                            sh "kubectl -n \${K8S_NAMESPACE} rollout restart deployment.apps/\${K8S_DEPLOYMENT_NAME}"
                        }
                    }
                }                
            }
        }
    }
    
    post { 
        always { 
            //sh "docker rmi test-order-test"
            cleanWs()
        }
    }
}