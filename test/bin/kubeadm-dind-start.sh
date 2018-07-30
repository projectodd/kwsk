#!/bin/bash -x

# Use a pod network that works on Travis
export POD_NETWORK_CIDR="10.244.0.0/16"

# Don't snapshot the cluster
export SKIP_SNAPSHOT=y

# Only run on 1 Node, for now
export NUM_NODES=1

# # Add the required cluster config for Knative serving
# export CONTROLLER_MANAGER_cluster_signing_cert_file="/var/lib/localkube/certs/ca.crt"
# export CONTROLLER_MANAGER_cluster_signing_key_file="/var/lib/localkube/certs/ca.key"
# export API_SERVER_admission_controller="DenyEscalatingExec,LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"

bash -x dind-cluster.sh clean
bash -x dind-cluster.sh up
