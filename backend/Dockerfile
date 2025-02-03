FROM golang:latest AS builder

WORKDIR /BackendApp

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /BackendApp/main .

CMD ["./main"]