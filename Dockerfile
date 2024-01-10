FROM golang:1.21-alpine3.18 AS builder

RUN apk add -U git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o s3-downloader

FROM alpine:latest

COPY --from=builder /build/s3-downloader /usr/bin/s3-downloader
