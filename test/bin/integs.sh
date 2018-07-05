#!/bin/bash

set -x

SCRIPTDIR=$(cd $(dirname "$0") && pwd)
TESTDIR="$SCRIPTDIR/.."

OWSK_HOME=$TESTDIR/openwhisk
KWSK_HOME=$1

if [ ! -d "$OWSK_HOME" ]; then
  git clone https://github.com/apache/incubator-openwhisk.git $OWSK_HOME
  sed -e "s:OPENWHISK_HOME:$OWSK_HOME:" <$TESTDIR/etc/whisk.properties >$OWSK_HOME/whisk.properties
fi
setsid go run ${KWSK_HOME}cmd/kwsk-server/main.go --port 8080 --istio $(minikube ip):32000 >kwsk.log 2>&1 &
KWSK_PID=$!

pushd $OWSK_HOME
./gradlew :tests:test --tests ${TESTS:-"system.basic.WskRest*"}
STATUS=$?
popd
kill -- "-$KWSK_PID"

exit $STATUS
