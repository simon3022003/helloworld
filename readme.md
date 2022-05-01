

# Introduction

This repo contains a small "Hello World" webserver which simulates a small microservice

## Functions


 - Create a docker image for the microservice. The smaller the image, the better.
 - Create all required resources in Kubernetes to expose the microservice to the public. 
 - Use MESSAGES env variable to configure the message displayed by the server
 - Create a K8S resource for scale up and down the microservice based on the CPU load
 - Create a Jenkins pipeline for deploying the microservice.
 - Expose the APIs for monitoring the K8s cluster.
