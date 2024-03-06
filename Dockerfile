FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/main cmd/go-api-test/main.go
CMD ["./bin/main"]
