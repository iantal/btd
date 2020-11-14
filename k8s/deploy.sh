#!/bin/bash -e

kubectl apply -f btd-service.yml
kubectl apply -f btd-deployment.yml