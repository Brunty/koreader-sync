# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /build/main ./cmd/koreader_sync_server
RUN go build -o /build/kor-cli ./cmd/koreader_sync

FROM alpine:latest

RUN apk add libc6-compat

WORKDIR /app

COPY --from=builder /build/main .
COPY --from=builder /build/kor-cli .

EXPOSE 8080

CMD ["./main"]