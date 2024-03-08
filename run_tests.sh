#!/bin/bash

docker-compose -f tests/docker-compose.test.yml up -d

sleep 10

go test -v ./tests/
TEST_STATUS=$?

docker-compose --profile test down -v

if [ $TEST_STATUS -eq 0 ]; then
  echo "✔️Tests was completed successfully."
else
  echo "✖️Tests was completed with errors."
  exit 1
fi
