#!/bin/bash

cd deployments

docker-compose -p crtest -f docker-compose.yaml -f docker-compose.test.yaml run integration-tests