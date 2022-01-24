FROM golang:1.16-alpine

ENV BIN_FILE /opt/cr/cr-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ${CODE_DIR}

RUN go build -ldflags "$LDFLAGS" -o ${BIN_FILE} ./internal/cmd

CMD ${BIN_FILE} update_stats
