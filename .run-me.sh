#!/bin/sh

echo 'Updating repo'
#apt-get update
echo 'Repo updated'


echo 'Installing docker, helm and minikube'
sudo apt-get install docker -y
sudo apt-get install minikube -y
sudo snap install helm

echo 'Adding bitnami to the repo'
helm repo add bitnami https://charts.bitnami.com/bitnami

echo 'Applying values to minikube'
helm install my-release bitnami/redis --values minikube_files/values-minikube.yml

echo 'Starting minikube'
minikube start

echo 'Applying kubectl configs'

kubectl apply -f app-configmap.yaml 
kubectl apply -f app-deployment.yaml 
kubectl apply -f app-secret.yaml


echo 'Environment ready'