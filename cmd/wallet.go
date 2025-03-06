package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/chtozamm/javacode-wallet/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	deposit  = "deposit"
	withdraw = "withdraw"
)

type operation struct {
	OperationType string `json:"operation_type"`
	Amount        int32  `json:"amount"`
}

func (app *application) handleOperation(w http.ResponseWriter, r *http.Request) {
	// Read and parse wallet UUID from path
	walletUUID := pgtype.UUID{}
	err := walletUUID.Scan(r.PathValue("wallet_id"))
	if err != nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	// Decode JSON from request to struct
	var op operation
	err = json.NewDecoder(r.Body).Decode(&op)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check operation type
	if op.OperationType != deposit && op.OperationType != withdraw {
		http.Error(w, "Unsupported operation type", http.StatusBadRequest)
		return
	}

	// Check amount
	if op.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Get current wallet balance
	oldBalance, err := app.queries.GetBalance(r.Context(), walletUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get wallet balance: %v\n", err)
		http.Error(w, "Failed to get wallet balance", http.StatusInternalServerError)
		return
	}

	// Calculate new balance
	var newBalance int32
	switch op.OperationType {
	case deposit:
		newBalance = oldBalance + op.Amount
	case withdraw:
		newBalance = oldBalance - op.Amount
	}

	// Check balance before withdrawal
	if op.OperationType == withdraw && newBalance < 0 {
		http.Error(w, fmt.Sprintf("Insufficient funds to withdraw: balance %d, trying to withdraw %d", oldBalance, op.Amount), http.StatusPaymentRequired)
		return
	}

	// Start transaction
	tx, err := app.db.Begin(r.Context())
	if err != nil {
		log.Printf("Failed to begin operation transaction: %v\n", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// Wrap queries with transaction
	queriesWithTx := app.queries.WithTx(tx)

	// Insert operation in database
	err = queriesWithTx.AddOperation(r.Context(), database.AddOperationParams{
		WalletID:      walletUUID,
		OperationType: op.OperationType,
		Amount:        op.Amount,
	})
	if err != nil {
		log.Printf("Failed to add operation: %v\n", err)
		http.Error(w, "Failed to add operation", http.StatusInternalServerError)
		return
	}

	// Update wallet balance
	updatedBalance, err := queriesWithTx.UpdateWallet(r.Context(), newBalance)
	if err != nil {
		http.Error(w, "Failed to update wallet balance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit transaction
	err = tx.Commit(r.Context())
	if err != nil {
		log.Printf("Failed to commit operation transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Write response with new balance
	w.Header().Set("Content-Type", "text/plain")
	_, err = fmt.Fprintf(w, "%d", updatedBalance)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (app *application) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	// Read and parse wallet UUID from path
	walletUUID := pgtype.UUID{}
	err := walletUUID.Scan(r.PathValue("wallet_id"))
	if err != nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	// Get current wallet balance
	balance, err := app.queries.GetBalance(r.Context(), walletUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get wallet balance: %v\n", err)
		http.Error(w, "Failed to get wallet balance", http.StatusInternalServerError)
		return
	}

	// Write response with current balance
	w.Header().Set("Content-Type", "text/plain")
	_, err = fmt.Fprintf(w, "%d", balance)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (app *application) handleGetWallets(w http.ResponseWriter, r *http.Request) {
	// Get wallets from the database
	wallets, err := app.queries.GetWallets(r.Context())
	if err != nil {
		log.Printf("Failed to get wallets: %v\n", err)
		http.Error(w, "Failed to retrieve wallets", http.StatusInternalServerError)
		return
	}

	// If no wallets are found, return an empty JSON array
	if len(wallets) == 0 {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("[]"))
		if err != nil {
			log.Printf("Failed to write response: %v\n", err)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
		return
	}

	// Marshal wallets slice into JSON
	walletsJSON, err := json.Marshal(wallets)
	if err != nil {
		log.Printf("Failed to marshal wallets into JSON: %v\n", err)
		http.Error(w, "Failed to marshal wallets", http.StatusInternalServerError)
		return
	}

	// Write response with wallets
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Explicitly set status to 200 OK
	_, err = w.Write(walletsJSON)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (app *application) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	// Create new empty wallet and return its ID
	walletID, err := app.queries.CreateWallet(r.Context(), 0)
	if err != nil {
		log.Printf("Failed to create wallet: %v\n", err)
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	// Write response with wallet ID
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(walletID.String()))
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (app *application) handleDeleteWallet(w http.ResponseWriter, r *http.Request) {
	// Read and parse wallet UUID from path
	walletUUID := pgtype.UUID{}
	err := walletUUID.Scan(r.PathValue("wallet_id"))
	if err != nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	// Delete the wallet
	err = app.queries.DeleteWallet(r.Context(), walletUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to delete wallet: %v\n", err)
		http.Error(w, "Failed to delete wallet", http.StatusInternalServerError)
		return
	}

	// Respond with no content status
	w.WriteHeader(http.StatusNoContent)
}
