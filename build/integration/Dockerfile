FROM golang:1.16-alpine

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0

CMD go get github.com/stretchr/testify/suite@v1.7.0
CMD go test -v -tags integration ./internal/tests/integration