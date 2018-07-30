#!/bin/bash

set -x

SCRIPTDIR=$(cd $(dirname "$0") && pwd)
TESTDIR="$SCRIPTDIR/.."

OWSK_HOME=$TESTDIR/openwhisk
KWSK_PORT=8180

if [ ! -d "$OWSK_HOME" ]; then
  git clone -b kwsk-tests --single-branch https://github.com/projectodd/incubator-openwhisk.git $OWSK_HOME
  cp $TESTDIR/etc/openwhisk-server-cert.pem $OWSK_HOME/ansible/roles/nginx/files/
fi
sed -e "s:OPENWHISK_HOME:$OWSK_HOME:; s/:8080/:$KWSK_PORT/" <$TESTDIR/etc/whisk.properties >$OWSK_HOME/whisk.properties

cat >$OWSK_HOME/tests/src/test/resources/application.conf <<EOF
# Blocking requests timeout by default after ~ 60s
akka.http.client.idle-timeout = 180 s
akka.http.host-connection-pool.idle-timeout = 180 s
akka.http.host-connection-pool.client.idle-timeout = 180 s
EOF


if [ "$(kubectl config current-context)" == "dind" ]; then
  ISTIO="localhost:32380"
else
  if [ "$(kubectl config current-context)" == "minikube" ]; then
    NODE_IP=$(minikube ip)
  else
    NODE_NAME=$(kubectl get node --no-headers | head -n 1 | awk '{print $1}')
    NODE_IP=$(kubectl get node ${NODE_NAME} -o 'jsonpath={.status.addresses[?(@.type=="InternalIP")].address}')
  fi
  NODE_PORT=$(kubectl get svc knative-ingressgateway -n istio-system -o 'jsonpath={.spec.ports[?(@.port==80)].nodePort}')
  ISTIO=${NODE_IP}:${NODE_PORT}
fi

nohup go run $TESTDIR/../cmd/kwsk-server/main.go --port $KWSK_PORT --istio $ISTIO --write-timeout 180s >kwsk.log 2>&1 &
KWSK_PID=$!

pushd $OWSK_HOME
TERM=dumb ./gradlew :tests:test --tests ${TESTS:-"system.basic.WskAction*"}
STATUS=$?
popd

pkill -P "$KWSK_PID"
exit $STATUS
