#!/bin/bash

cd deployments

DB_PORT=8086: RABBIT_PORT=8087: docker compose -p cr-database-rabbit -f docker-compose.yaml up db migrations rabbit