# Build Binaries in Alpine
# =============================================================================

FROM golang:1.12-alpine AS builder
WORKDIR /app

RUN apk update --no-cache
RUN apk add --no-cache git bash make ca-certificates git gcc g++ libc-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=arm go build -o ./bin/rover ./cmd/rover/main.go 


# Copy Binaries to Raspbian Server
# =============================================================================

FROM balenalib/rpi-raspbian AS server
WORKDIR /app

RUN apt-get -q update && apt-get -y install libraspberrypi-bin
RUN apt-get clean && rm -rf /var/lib/apt/lists/*
RUN usermod -a -G video root

COPY --from=builder /app .

CMD ["/app/bin/rover"]
