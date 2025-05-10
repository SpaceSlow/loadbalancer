FROM golang:1.23-alpine

COPY . /app
WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o loadbalancer cmd/loadbalancer/main.go
