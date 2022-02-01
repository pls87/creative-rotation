# Сreative Rotation
## Run
```bash
make run-docker-api-with-tool #runs components in containers using docker-compose
```
## Test
```bash
make run-docker-integration-test #runs integration tests using docker containers
```
## Components
### Storage
PostgreSQL database
### API
REST API available here: http://127.0.0.1:8080/cr
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
Providing simple UI for API access. Available here: http://127.0.0.1:8080/swaggerui
### Stats Updater
Background job to update impression/conversion aggregated statistics
### Adminer
Lightweight database admin tool. Available here: http://127.0.0.1:8080/dbadmin
## Configuration
Configuration via both environment variables and config file is supported. Sample configuration file is [here](https://github.com/pls87/creative-rotation/blob/develop/configs/sample.toml)
