#!/bin/bash

API_URL=http://localhost:8080/api/v1
DB_URL=postgres://javacode:secret@localhost:5432/wallet

case "$1" in
    start)
		docker compose up -d && goose postgres -dir sql/schema $DB_URL up
        ;;
    stop)
		docker compose down
        ;;
    bench)
        # Create new wallet
        wallet_id=$(curl -s $API_URL/wallets -X POST)
		# Run benchmark
        ab -n 10000 -c 1000 $API_URL/wallets/$wallet_id
        # Delete wallet
        curl -s $API_URL/wallets/$wallet_id -X DELETE         
        ;;
    clean)
		docker compose down && docker image rm wallet-server && docker volume rm javacode-wallet_postgres_data
        ;;
    *)
        echo "Usage: $0 {start|stop|bench|clean}"
        exit 1
esac
