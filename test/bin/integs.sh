#!/bin/bash

set -x

SCRIPTDIR=$(cd $(dirname "$0") && pwd)
TESTDIR="$SCRIPTDIR/.."

OWSK_HOME=$TESTDIR/openwhisk
KWSK_PORT=8181

if [ ! -d "$OWSK_HOME" ]; then
  git clone -b kwsk-tests --single-branch https://github.com/projectodd/incubator-openwhisk.git $OWSK_HOME
  cp $TESTDIR/etc/openwhisk-server-cert.pem $OWSK_HOME/ansible/roles/nginx/files/
fi
sed -e "s:OPENWHISK_HOME:$OWSK_HOME:; s/:8080/:$KWSK_PORT/" <$TESTDIR/etc/whisk.properties >$OWSK_HOME/whisk.properties

ISTIO="localhost:8181"
nohup kubectl port-forward -n istio-system deployment/knative-ingressgateway 8880:80 >portforward.log 2>&1 &
PORTFORWARD_PID=$!

nohup go run $TESTDIR/../cmd/kwsk-server/main.go --port $KWSK_PORT --istio $ISTIO >kwsk.log 2>&1 &
KWSK_PID=$!

pushd $OWSK_HOME
TERM=dumb ./gradlew :tests:test --tests ${TESTS:-"system.basic.WskRest*"}
STATUS=$?
popd

pkill -P "$KWSK_PID"
pkill -P "$PORTFORWARD_PID"

exit $STATUS
