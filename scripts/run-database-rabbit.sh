#!/bin/sh

cd deployments

DB_PORT=8086 RABBIT_PORT=5673 docker compose -p cr-database-rabbit -f docker-compose.yaml up db migrations rabbit