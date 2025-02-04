FROM golang:latest AS builder

WORKDIR /PingerApp

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o pinger .

FROM alpine:latest

RUN apk add --no-cache libc6-compat
RUN apk add --no-cache docker-cli

WORKDIR /root/

COPY --from=builder /PingerApp/pinger .
COPY --from=builder /PingerApp/ping /root/ping


CMD ["./pinger"]
