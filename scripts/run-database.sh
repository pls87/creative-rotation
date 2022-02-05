#!/bin/sh

cd deployments

DB_PORT=8086 docker compose -p cr-database -f docker-compose.yaml up db migrations