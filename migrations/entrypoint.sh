#!/bin/bash

DBSTRING="host=$POSTGRES_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable"

for i in $(seq 1 5); do
    goose postgres "$DBSTRING" up && break
    sleep 2
done
