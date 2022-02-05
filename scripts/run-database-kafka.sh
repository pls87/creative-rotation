#!/bin/sh

cd deployments

DB_PORT=8086 KAFKA_PORT=9094 docker compose -p cr-database-kafka -f docker-compose.yaml up db migrations broker