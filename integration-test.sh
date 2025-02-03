#!/bin/bash

set -e
set -o pipefail

START_TIME=$(date +%s)
SERVER_CONTAINER="server"
NUM_CLIENTS=${1:-10}

make clean-docker
make network

echo "Starting the server..."
make run-server

echo "Waiting for the server to start..."
sleep 1

docker logs -f "$SERVER_CONTAINER" | tee server.log & LOG_PID=$!

echo "Starting $NUM_CLIENTS clients in parallel..."

EXIT_CODE=0
PIDS=()
for i in $(seq 1 "$NUM_CLIENTS"); do
    (
        echo "🕒 Running client $i out of $NUM_CLIENTS..."
        if timeout 30s make run-client; then
            echo "✅ Client $i out of $NUM_CLIENTS passed"
        else
            echo "❌ Client $i out of $NUM_CLIENTS failed" >&2
            EXIT_CODE=1
        fi
    ) & PIDS+=($!)
done

echo "Waiting for all clients to finish..."
for pid in "${PIDS[@]}"; do
    wait "$pid" || EXIT_CODE=1
done

kill "$LOG_PID" || true

echo "Cleaning up docker resources"
make clean-docker

END_TIME=$(date +%s)
ELAPSED_TIME=$((END_TIME - START_TIME))
echo "Total execution time: ${ELAPSED_TIME} seconds"

exit $EXIT_CODE
