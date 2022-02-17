# Сreative Rotation

## Run

To run entire application with additional tools like database adminer use the command below

```bash
make run #runs components in containers using docker-compose
```

To run api and status_updater processes locally you can use the command below

```bash
make run-local #runs api and stats_updater processes on host machine
```

You can run both application components (API and Stats Updater) separately by using commands below

```bash
make run-api-local #runs http api on the host machine
make run-stats-updater-local #runs stats updater on the host machine
```

It's recommended to run Postgres and rabbitMQ in docker using the following command before launching api and
status_updater locally

```bash
make run-database-rabbit #runs api and stats_updater processes on host machine
```

## Test

```bash
make test-unit #runs unit tests
```

```bash
make run-integration-test #runs integration tests using docker containers
```

```bash
make test #runs all tests
```

## Components

### Storage

PostgreSQL database

### API

If you run entire application in docker then REST API available here: http://127.0.0.1:8080/cr

#### Endpoints

```
GET     "/creative" - get list of all creatives
POST    "/creative" - create new creative
POST    "/creative/{creative_id:[0-9]+}/slot" - attach creative to slot
DELETE  "/creative/{creative_id:[0-9]+}/slot/{slot_id:[0-9]+}" - detach creative out of slot
POST    "/conversion" - track click on creative
POST    "/impression" - track impression for creative
GET     "/creative/next" - get next creative to show
GET     "/slot" - get list of all slots
POST    "/slot" - create new slot
GET     "/segment" - get list of all segments
POST    "/segment" - create new segment
```

### Swagger UI

Provides simple UI for API access. If you run entire application in docker then swagger is available
here: http://127.0.0.1:8080/swaggerui

### Stats Updater

Background job to update impression/conversion aggregated statistics

### Adminer

Lightweight database admin tool. If you run entire application in docker then adminer is available
here: http://127.0.0.1:8080/dbadmin

## Configuration

Configuration via both environment variables and config file is supported. Sample configuration file
is [here](https://github.com/pls87/creative-rotation/blob/develop/configs/sample.toml)
