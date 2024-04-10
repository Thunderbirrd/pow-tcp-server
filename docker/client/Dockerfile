FROM golang:alpine as go

ENV BUILD_DIR /src

WORKDIR /${BUILD_DIR}

COPY go.mod go.sum ./
RUN  go mod download

COPY . ./
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 \
    go build -o app ./cmd/client/main.go

FROM alpine:latest

ENV BUILD_DIR /src

COPY --from=go ${BUILD_DIR} ${BUILD_DIR}

CMD ${BUILD_DIR}/app