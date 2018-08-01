#!/bin/bash
set -x

# install istio
curl -L https://storage.googleapis.com/knative-releases/serving/latest/istio.yaml \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply -f -

# Don't try to inject in the istio-system namespace
kubectl label namespace istio-system istio-injection=disabled
# label the default namespace with istio-injection=enabled.
kubectl label namespace default istio-injection=enabled

TIMEOUT=600 # in seconds
INTERVAL=5  # in seconds

ELAPSED=0
PASSED=false
sleep $INTERVAL
until [ $ELAPSED -ge $TIMEOUT ]; do
  kubectl get pods -n istio-system | grep -v -E "(Running|Completed|STATUS)"
  EXIT_STATUS=$?
  if [ $EXIT_STATUS -eq 1 ]; then
    PASSED=true
    break
  fi

  let ELAPSED=ELAPSED+$INTERVAL
  sleep $INTERVAL
done
if [ "$PASSED" = false ]; then
  echo "Failed to deploy Istio within $TIMEOUT seconds"
  kubectl get pods -o wide -n istio-system
  kubectl describe pods -n istio-system
  kubectl get all --all-namespaces
  kubectl describe node
  kubectl get events
  exit 1
fi


# install knative
curl -L https://storage.googleapis.com/knative-releases/serving/latest/release-no-mon.yaml \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply -f -

ELAPSED=0
PASSED=false
sleep $INTERVAL
until [ $ELAPSED -ge $TIMEOUT ]; do
  kubectl get pods -n knative-serving | grep -v -E "(Running|Completed|STATUS)"
  EXIT_STATUS=$?
  if [ $EXIT_STATUS -eq 1 ]; then
    PASSED=true
    break
  fi

  let ELAPSED=ELAPSED+$INTERVAL
  sleep $INTERVAL
done
if [ "$PASSED" = false ]; then
  echo "Failed to deploy Knative within $TIMEOUT seconds"
  kubectl get pods -o wide -n knative-serving
  kubectl describe pods -n knative-serving
  kubectl get all --all-namespaces
  kubectl describe node
  kubectl get events
  exit 1
fi

# install knative eventing
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release.yaml
# and the stub bus
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release-bus-stub.yaml
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release-clusterbus-stub.yaml
# and the Kubernetes Events source
kubectl apply -f https://storage.googleapis.com/knative-releases/eventing/latest/release-source-k8sevents.yaml

set +x
echo "Knative successfully installed!"
