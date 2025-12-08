FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest

RUN apk add --no-cache docker-cli

WORKDIR /root/
COPY --from=builder /app/main .

ENTRYPOINT ["./main"]