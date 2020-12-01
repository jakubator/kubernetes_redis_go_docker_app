#!/bin/sh

echo 'Updating repo'
#apt-get update
echo 'Repo updated'


echo 'Installing docker, helm and minikube'
sudo apt-get install docker -y
sudo apt-get install curl -y
sudo apt install virtualbox virtualbox-ext-pack -y
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
sudo chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl

curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
sudo dpkg -i minikube_latest_amd64.deb
sudo snap install helm --classic

echo 'Adding bitnami to the repo'
helm repo add bitnami https://charts.bitnami.com/bitnami

echo 'Applying values to minikube'
helm install my-release bitnami/redis --values values-minikube.yml

echo 'Starting minikube'
minikube start --driver=virtualbox

echo 'Applying kubectl configs'

kubectl apply -f app-configmap.yaml 
kubectl apply -f app-deployment.yaml 
kubectl apply -f app-secret.yaml


echo 'Environment ready'
