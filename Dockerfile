# https://container-solutions.com/faster-builds-in-docker-with-go-1-11/

FROM golang:1.12-alpine as build_base
WORKDIR /app

RUN apk update --no-cache
RUN apk add --no-cache git bash make ca-certificates git gcc g++ libc-dev

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM build_base AS server_builder
