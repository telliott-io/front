#!/bin/sh
#kubectl create configmap kubeenv --from-literal=context=`kubectl config current-context` --namespace=front

echo `kubectl config current-context` > kubecontext.txt