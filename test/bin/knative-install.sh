#!/bin/bash
set -x

# install istio
wget -O - https://storage.googleapis.com/knative-releases/latest/istio.yaml \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply -f -
# label the default namespace with istio-injection=enabled.
kubectl label namespace default istio-injection=enabled
# wait until each istio pod is up
while kubectl get pods -n istio-system | grep -v -E "(Running|Completed|STATUS)"; do
    sleep 5
done

# install knative
kubectl apply -f https://storage.googleapis.com/knative-releases/latest/release-lite.yaml
# wait until each knative pod is up
while kubectl get pods -n knative-serving | grep -v -E "(Running|Completed|STATUS)"; do
    sleep 5
done

set +x
echo "Knative successfully installed!"
