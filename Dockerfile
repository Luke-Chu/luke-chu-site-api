# syntax=docker/dockerfile:1.7

FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/luke-chu-site-api ./cmd/server

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S -G app app && \
    update-ca-certificates

COPY --from=builder /out/luke-chu-site-api /app/luke-chu-site-api
COPY configs /app/configs

USER app

ENV APP_ENV=prod \
    GIN_MODE=release \
    TZ=Asia/Shanghai

EXPOSE 8080

ENTRYPOINT ["/app/luke-chu-site-api"]
