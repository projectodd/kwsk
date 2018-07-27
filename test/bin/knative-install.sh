#!/bin/bash
set -x

# install istio
curl -L https://storage.googleapis.com/knative-releases/serving/latest/istio.yaml \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply -f -

# label the default namespace with istio-injection=enabled.
kubectl label namespace default istio-injection=enabled

TIMEOUT=600 # in seconds
INTERVAL=5  # in seconds

ELAPSED=0
PASSED=false
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
  kubectl describe node
  kubectl get events
  exit 1
fi

set +x
echo "Knative successfully installed!"
