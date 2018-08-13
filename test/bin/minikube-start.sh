#!/bin/bash

K8S_VERSION=${K8S_VERSION:-v1.10.5}
VM_DRIVER=${VM_DRIVER:-kvm2}
BOOTSTRAPPER=${BOOTSTRAPPER:-kubeadm}
MINIKUBE_MEMORY=${MINIKUBE_MEMORY:-8192}
MINIKUBE_CPUS=${MINIKUBE_CPUS:-4}

if [ "$CI" = "true" ]; then
  K8S_VERSION=v1.10.0
  VM_DRIVER=none
  BOOTSTRAPPER=localkube
fi

exec minikube start --memory=$MINIKUBE_MEMORY --cpus=$MINIKUBE_CPUS \
     --kubernetes-version=$K8S_VERSION \
     --vm-driver=$VM_DRIVER \
     --bootstrapper=$BOOTSTRAPPER \
     --extra-config=controller-manager.cluster-signing-cert-file="/var/lib/localkube/certs/ca.crt" \
     --extra-config=controller-manager.cluster-signing-key-file="/var/lib/localkube/certs/ca.key" \
     --extra-config=apiserver.admission-control="DenyEscalatingExec,LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"
