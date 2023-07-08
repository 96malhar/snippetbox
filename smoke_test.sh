#!/bin/bash

cleanup() {
  # Stop the web server if it is still running
  kill "$PID"

  # Clean up the binary
  rm -rf ./bin
}

# Build the Go binary
go build -o ./bin/server ./cmd/web

# Start the web server in the background
./bin/server &
PID=$!

# Define the cleanup function to be executed on script exit
trap cleanup EXIT ERR

# Wait for the server to start
sleep 2

# Make a GET request to the server and check the response
RESPONSE=$(curl -k -s https://localhost:4000/ping)

# Check if the response contains the expected string
EXPECTED="OK"
if [[ $RESPONSE == "$EXPECTED" ]]; then
  echo "Smoke test passed. Server is running correctly."
  exit 0
else
  echo "Smoke test failed. Server did not respond as expected."
  echo "Expected: $EXPECTED"
  echo "Actual  : $RESPONSE"
  exit 1
fi
