#!/bin/bash

docker-compose --profile test up -d

sleep 7

go test -v ./tests/
TEST_STATUS=$?

docker-compose --profile test down -v

if [ $TEST_STATUS -eq 0 ]; then
  echo "✔️Tests was completed successfully."
else
  echo "✖️Tests was completed with errors."
  exit 1
fi
