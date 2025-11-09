FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache git bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o app main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

COPY config.json .

EXPOSE 9000

CMD ["./app"]
