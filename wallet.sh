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
    test)
        # Create new wallet
        wallet_id=$(curl -s $API_URL/wallets -X POST)
        if [ $? -ne 0 ]; then
            echo "Error creating wallet"
            exit 1
        fi
        echo "Your wallet_id: $wallet_id"

        # Get balance
        balance=$(curl -s $API_URL/wallets/$wallet_id)
        if [ $? -ne 0 ]; then
            echo "Error getting balance"
            exit 1
        fi
        echo "Balance: $balance"

        # Deposit 
        curl -s $API_URL/wallets/$wallet_id -X POST \
            -H "Content-Type: application/json" \
            -d '{"operation_type": "deposit","amount": 500}'
        if [ $? -ne 0 ]; then
            echo "Error depositing funds"
            exit 1
        fi

        # Get balance
        balance=$(curl -s $API_URL/wallets/$wallet_id)
        if [ $? -ne 0 ]; then
            echo "Error getting balance after deposit"
            exit 1
        fi
        echo "Balance after deposit: $balance"

        # Withdraw 
        curl -s $API_URL/wallets/$wallet_id -X POST \
            -H "Content-Type: application/json" \
            -d '{"operation_type": "withdraw","amount": 150}'
        if [ $? -ne 0 ]; then
            echo "Error withdrawing funds"
            exit 1
        fi

        # Get balance
        balance=$(curl -s $API_URL/wallets/$wallet_id)
        if [ $? -ne 0 ]; then
            echo "Error getting balance after withdrawal"
            exit 1
        fi
        echo "Balance after withdrawal: $balance"

        # Try to withdraw more than the balance holds
        message=$(curl -s $API_URL/wallets/$wallet_id -X POST \
            -H "Content-Type: application/json" \
            -d '{"operation_type": "withdraw","amount": 10000}')
        if [ $? -ne 0 ]; then
            echo "Error withdrawing funds"
            exit 1
        fi
        echo "Response when trying to withdraw more than the wallet has: $message"

        # Get all wallets
        all_wallets=$(curl -s $API_URL/wallets -u javacode:secret)
        if [ $? -ne 0 ]; then
            echo "Error getting all wallets"
            exit 1
        fi
        echo "All wallets: $all_wallets"

        # Delete wallet
        curl -s $API_URL/wallets/$wallet_id -X DELETE
        if [ $? -ne 0 ]; then
            echo "Error deleting wallet"
            exit 1
        fi
        echo "Wallet deleted successfully"

        # Get all wallets after deletion
        all_wallets=$(curl -s $API_URL/wallets -u javacode:secret)
        if [ $? -ne 0 ]; then
            echo "Error getting all wallets"
            exit 1
        fi
        echo "All wallets after deletion: $all_wallets"
    ;;

    clean)
		docker compose down && docker image rm wallet-server && docker volume rm javacode-wallet_postgres_data
        ;;
    *)
        echo "Usage: $0 {start|stop|bench|test|clean}"
        exit 1
esac
