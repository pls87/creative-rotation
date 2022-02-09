FROM golang:1.16-alpine

RUN apk update && \
    apk upgrade && \
    apk add postgresql-client && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0

RUN go get github.com/stretchr/testify/suite@v1.7.0
CMD go test -v -tags integration ./internal/tests/integration