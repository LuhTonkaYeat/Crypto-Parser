FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /server api/server.go

RUN go build -o /cli cmd/cli/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /server /app/server
COPY --from=builder /cli /app/cli

EXPOSE 8080

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]