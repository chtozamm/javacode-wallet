FROM docker.io/golang:1.24 AS builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o wallet-server -ldflags="-s -w" ./cmd 

FROM docker.io/alpine:latest
RUN mkdir /app && adduser -h /app -D javacode
WORKDIR /app
COPY --chown=javacode --from=builder /app/.env .
COPY --chown=javacode --from=builder /app/wallet-server .
EXPOSE 8080

CMD ["./wallet-server"]
