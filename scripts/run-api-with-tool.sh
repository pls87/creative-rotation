#!/bin/sh

cd deployments

docker compose -p cr -f docker-compose.yaml -f docker-compose.tool.yaml  up