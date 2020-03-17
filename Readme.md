# Covid 19 updates

This repository houses a go program used to retrieve statistics of the COVID19 Virus.

There are two folders in this project.

### Lambda

The lambda folder includes a go program written for AWS Lambda that will post messages to a `slack` webhook. Instructions for creating that are found below.

-	build go binary with the following env to convert it for AWS

`env GOOS=linux GOARCH=amd64 go build -o covid`

-   Zip the file with

`zip -j covid.zip covid`

-   Upload to lambda and name the handler "covid"


### k8s
The k8s folder includes a go program and a Dockerfile for creating a container image. There is also a cronjob created for use in a Kubernetes cluster.
