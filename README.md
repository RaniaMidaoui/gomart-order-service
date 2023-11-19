# goMart-order-service
This repository contains the code and Dockerfile for the order microservice of the **goMart-commerce** application, along with the Jenkinsfile describing the CI/CD pipeline for the microservice.

To run the code, you need to have Golang package installed:

1- Download the package from [the official website](https://go.dev/doc/install)
##### For Linux users:
2- Remove any previous Go installation  then extract the archive you just downloaded into /usr/local:
```
 $ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
```
3- Add /usr/local/go/bin to the PATH environment variable:
```
$ export PATH=$PATH:/usr/local/go/bin
```
4- Verify that you've installed Go:
```
$ go version
```
##### For Windows users:
Follow the prompt after opening the MSI file you downloaded from [the official website](https://go.dev/doc/install).

### To test the microservice locally:
Start by making sure all the dependencies are installed and run the code, it will tell you that it's listenning on the application port configured:
```
$ go mod tidy
$ make proto
$ make server
```
To test the microservice, the [API Gateway](https://github.com/RaniaMidaoui/goMart-gateway) must be running in order to redirect the request to the order microservice, you must already have registered and logged in a user with the [authentication microservice](https://github.com/RaniaMidaoui/goMart-authentication-service) and got his authorization token (\$TOKEN) and created a product with the [product microservice](https://github.com/RaniaMidaoui/goMart-product-service) and got its ID (\$PRODUCT_ID):
```
#Create order
curl --request POST \
  --url http://localhost:3000/order \
  --header "Authorization: Bearer $TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
 "productId": $(PRODUCT_ID),
 "quantity": 1
}'
```

