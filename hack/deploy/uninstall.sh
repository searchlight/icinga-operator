#!/bin/bash
set -x

kubectl delete deployment -l app=searchlight -n kube-system
kubectl delete service -l app=searchlight -n kube-system
kubectl delete secret -l app=searchlight -n kube-system

# Delete RBAC objects, if --rbac flag was used.
kubectl delete serviceaccount -l app=searchlight -n kube-system
kubectl delete clusterrolebindings -l app=searchlight -n kube-system
kubectl delete clusterrole -l app=searchlight -n kube-system
