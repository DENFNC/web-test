# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o websrv ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/websrv .
EXPOSE 8080
CMD ["./websrv"] 