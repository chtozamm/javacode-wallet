FROM docker.io/golang:1.24 AS builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -C cmd -o wallet-server -ldflags="-s -w" . 

FROM docker.io/alpine:latest
RUN mkdir /app && adduser -h /app -D wallet-server
WORKDIR /app
COPY --chown=wallet-server --from=builder /app/.env .
COPY --chown=wallet-server --from=builder /app/cmd/wallet-server .
EXPOSE 8080

CMD ["./wallet-server"]
