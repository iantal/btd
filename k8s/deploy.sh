#!/bin/bash -e

kubectl apply -f btd-configmap.yml
kubectl apply -f btd-pv.yml
kubectl apply -f btd-pvc.yml
kubectl apply -f btd-service.yml
kubectl apply -f btd-statefulset.yml