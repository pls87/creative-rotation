#!/bin/bash

cd deployments

trap "docker compose -p crtest -f docker-compose.yaml -f docker-compose.tests.yaml down" EXIT

docker compose -p crtest -f docker-compose.yaml -f docker-compose.test.yaml run integration-tests