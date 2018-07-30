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
ISTIO=$(minikube ip):$(kubectl get svc knative-ingressgateway -n istio-system -o 'jsonpath={.spec.ports[?(@.port==80)].nodePort}')

nohup go run $TESTDIR/../cmd/kwsk-server/main.go --port $KWSK_PORT --istio $ISTIO >kwsk.log 2>&1 &
KWSK_PID=$!

pushd $OWSK_HOME
TERM=dumb ./gradlew :tests:test --tests ${TESTS:-"system.basic.WskRest*"}
STATUS=$?
popd

pkill -P "$KWSK_PID"
exit $STATUS
