#!/bin/bash
set -x

# One of release.yaml, release-lite.yaml, or release-no-mon.yaml
KNATIVE_SERVING_FLAVOR=${KNATIVE_SERVING_FLAVOR:-release-no-mon.yaml}

wait_for_pods() {
  namespace=$1
  timeout=600 # in seconds
  interval=5  # in seconds
  elapsed=0
  passed=false
  sleep $interval
  until [ $elapsed -ge $timeout ]; do
    kubectl get pods -n $namespace | grep -v -E "(Running|Completed|STATUS)"
    exit_status=$?
    if [ $exit_status -eq 1 ]; then
      passed=true
      break
    fi

    let elapsed=elapsed+$interval
    sleep $interval
  done
  if [ "$passed" = false ]; then
    echo "Failed to deploy pods in $namespace within $timeout seconds"
    kubectl get pods -o wide -n $namespace
    kubectl describe pods -n $namespace
    kubectl get all --all-namespaces
    kubectl describe node
    kubectl get events
    exit 1
  fi
}

# install istio
curl -L https://storage.googleapis.com/knative-releases/serving/latest/istio.yaml \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply -f -

# Don't try to inject in the istio-system namespace
kubectl label namespace istio-system istio-injection=disabled
# label the default namespace with istio-injection=enabled.
kubectl label namespace default istio-injection=enabled

wait_for_pods "istio-system"


# install knative
curl -L https://storage.googleapis.com/knative-releases/serving/latest/${KNATIVE_SERVING_FLAVOR} \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply -f -

wait_for_pods "knative-serving"

# install knative eventing
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release.yaml
wait_for_pods "knative-eventing"

# and the stub bus
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release-bus-stub.yaml
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release-clusterbus-stub.yaml
# and the Kubernetes Events source
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release-source-k8sevents.yaml

set +x
echo "Knative successfully installed!"
