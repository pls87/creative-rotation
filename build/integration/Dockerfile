FROM golang:1.16-alpine

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0

CMD go test -v -tags integration ./internal/tests/integration