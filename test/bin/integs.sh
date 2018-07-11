#!/bin/bash

set -x

SCRIPTDIR=$(cd $(dirname "$0") && pwd)
TESTDIR="$SCRIPTDIR/.."

OWSK_HOME=$TESTDIR/openwhisk

if [ ! -d "$OWSK_HOME" ]; then
  git clone https://github.com/apache/incubator-openwhisk.git $OWSK_HOME
  cp $TESTDIR/etc/openwhisk-server-cert.pem $OWSK_HOME/ansible/roles/nginx/files/
fi
sed -e "s:OPENWHISK_HOME:$OWSK_HOME:" <$TESTDIR/etc/whisk.properties >$OWSK_HOME/whisk.properties
ISTIO=$(minikube ip):$(kubectl get svc knative-ingressgateway -n istio-system -o 'jsonpath={.spec.ports[?(@.port==80)].nodePort}')
setsid go run cmd/kwsk-server/main.go --port 8080 --istio $ISTIO >kwsk.log 2>&1 &
KWSK_PID=$!

pushd $OWSK_HOME
./gradlew :tests:test --tests ${TESTS:-"system.basic.WskRest*"}
STATUS=$?
popd
kill -- "-$KWSK_PID"

exit $STATUS
