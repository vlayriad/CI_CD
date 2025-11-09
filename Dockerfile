# 1. Build stage
FROM golang:1.24.6 AS builder

WORKDIR /app

# copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# copy all source code
COPY . .

# build the binary
RUN go build -v -o app main.go

# 2. Final stage
FROM alpine:latest

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/app .

# copy config.json / env if needed
COPY config.json .

# expose port
EXPOSE 9000

# run the binary
CMD ["./app"]
