# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /build/main .

WORKDIR /app/cmd/koreader_sync

RUN go build -o /build/kor-cli .

FROM alpine:latest

RUN apk add libc6-compat

WORKDIR /app

COPY --from=builder /build/main .
COPY --from=builder /build/kor-cli .

EXPOSE 8080

CMD ["./main"]