#!/bin/sh
##
# Script to deploy a Kubernetes project with a StatefulSet running a MongoDB Replica Set, to a local Minikube environment.
##

# Create keyfile for the MongoD cluster as a Kubernetes shared secret
TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 > $TMPFILE
kubectl create secret generic shared-bootstrap-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE --from-file=my-pem=../ssl/mongodb.pem
rm $TMPFILE

# Create mongodb service with mongod stateful-set
kubectl apply -f ../deployment/deployment.yaml
sleep 5

# Print current deployment state (unlikely to be finished yet)
kubectl get all 
kubectl get persistentvolumes
echo
echo "Keep running the following command until all 'mongod-n' pods are shown as running:  kubectl get all"
echo

