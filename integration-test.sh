#!/bin/bash

set -e
set -o pipefail

SERVER_CONTAINER="server"

make clean-docker

echo "Starting the server..."
make run-server

echo "Waiting for the server to start..."
sleep 1

docker logs -f "$SERVER_CONTAINER" & LOG_PID=$!

echo "Starting the client..."
if make run-client; then
    echo "✅ Test passed: make run-client exited with code 0"
    EXIT_CODE=0
else
    echo "❌ Error: make run-client exited with a non-zero code" >&2
    EXIT_CODE=1
fi

kill "$LOG_PID" || true

echo "Cleaning up docker resources"
make clean-docker

exit $EXIT_CODE
