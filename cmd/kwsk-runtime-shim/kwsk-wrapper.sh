#!/bin/bash

# Start the runtime shim
echo "Starting shim"
kwsk-runtime-shim &
status=$?
echo "Shim Status: ${status}"
shim_pid=$!
echo "Shim PID: ${shim_pid}"
if [ $status -ne 0 ]; then
  echo "Failed to start runtime shim"
  exit $status
fi

# Start the regular server
echo "Starting server: $@"
"$@" &
status=$?
echo "Server Status: ${status}"
server_pid=$!
echo "Server PID: ${server_pid}"
if [ $status -ne 0 ]; then
  echo "Failed to start server"
  exit $status
fi

while sleep 10; do
  ps | grep ${shim_pid} | grep -q -v grep
  shim_status=$?
  ps | grep ${server_pid} | grep -q -v grep
  server_status=$?
  if [ $shim_status -ne 0 -o $server_status -ne 0 ]; then
    echo "ERROR: shim or server process has terminated."
    kill $shim_pid
    kill $server_pid
    exit 1
  fi
done
