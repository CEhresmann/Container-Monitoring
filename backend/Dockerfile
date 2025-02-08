FROM golang:latest AS builder

WORKDIR /BackendApp

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o backend .

FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /BackendApp/backend .
COPY --from=builder /BackendApp/config /root/config
COPY --from=builder /BackendApp/db /root/db
COPY --from=builder /BackendApp/my_handler /root/my_handler
COPY --from=builder /BackendApp/queue /root/queue




CMD ["./backend"]