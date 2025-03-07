package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chtozamm/javacode-wallet/internal/database"
	"github.com/chtozamm/javacode-wallet/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetBalance(t *testing.T) {
	validUUID := "fe6403a7-8b42-4449-abe6-a8508199a0d4"
	invalidUUID := "fe6403a7-8b421-449-abe6-a8508199a0d4"

	tests := []struct {
		name         string
		walletID     string
		mockBalance  int32
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid wallet ID",
			walletID:     validUUID,
			mockBalance:  100,
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: "100\n",
		},
		{
			name:         "Wallet not found",
			walletID:     validUUID,
			mockBalance:  0,
			mockError:    sql.ErrNoRows,
			expectedCode: http.StatusNotFound,
			expectedBody: "Wallet not found\n",
		},
		{
			name:         "Invalid wallet ID",
			walletID:     invalidUUID,
			mockBalance:  0,
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid wallet ID\n",
		},
		{
			name:         "Unexpected error",
			walletID:     validUUID,
			mockBalance:  100,
			mockError:    errors.New("Unexpected error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to get wallet balance\n",
		},
	}

	for _, tc := range tests {
		mockDB := &mocks.DBTX{
			Balance: tc.mockBalance,
			Err:     tc.mockError,
		}

		t.Run(tc.name, func(t *testing.T) {
			// Create a new application with the mock queries
			app := &application{
				queries: database.New(mockDB),
			}

			// Create and configure a new HTTP request
			req := httptest.NewRequest("GET", "/api/v1/wallets/"+tc.walletID, nil)
			req.SetPathValue("wallet_id", tc.walletID)

			// Create a new response recorder
			w := httptest.NewRecorder()

			// Call the handler
			app.handleGetBalance(w, req)

			// Check the response
			res := w.Result()
			assert.Equal(t, tc.expectedCode, res.StatusCode)

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}

func TestHandleOperation(t *testing.T) {
	validUUID := "fe6403a7-8b42-4449-abe6-a8508199a0d4"
	invalidUUID := "fe6403a7-8b421-449-abe6-a8508199a0d4"

	tests := []struct {
		name         string
		walletID     string
		mockBalance  int32
		mockError    error
		op           operation
		expectedCode int
		expectedBody string
	}{
		// NOTE: the commented test cases below cause panic due to the unimplemented database connection for transaction mocking...
		// Using two tables and transactions has made testing more complex.
		// {
		// 	name:         "Valid deposit operation",
		// 	walletID:     validUUID,
		// 	mockBalance:  100,
		// 	mockError:    nil,
		// 	op:           operation{OperationType: deposit, Amount: 50},
		// 	expectedCode: http.StatusNoContent,
		// 	expectedBody: "",
		// },
		// {
		// 	name:         "Valid withdraw operation",
		// 	walletID:     validUUID,
		// 	mockBalance:  100,
		// 	mockError:    nil,
		// 	op:           operation{OperationType: withdraw, Amount: 30},
		// 	expectedCode: http.StatusNoContent,
		// 	expectedBody: "",
		// },
		{
			name:         "Invalid wallet ID",
			walletID:     invalidUUID,
			mockBalance:  0,
			mockError:    nil,
			op:           operation{OperationType: withdraw, Amount: 100},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid wallet ID\n",
		},
		{
			name:         "Insufficient funds for withdrawal",
			walletID:     validUUID,
			mockBalance:  50,
			mockError:    nil,
			op:           operation{OperationType: withdraw, Amount: 100},
			expectedCode: http.StatusPaymentRequired,
			expectedBody: "Insufficient funds to withdraw: balance 50, trying to withdraw 100\n",
		},
		{
			name:         "Wallet not found",
			walletID:     validUUID,
			mockBalance:  0,
			mockError:    sql.ErrNoRows,
			op:           operation{OperationType: deposit, Amount: 50},
			expectedCode: http.StatusNotFound,
			expectedBody: "Wallet not found\n",
		},
		{
			name:         "Invalid operation type",
			walletID:     validUUID,
			mockBalance:  100,
			mockError:    nil,
			op:           operation{OperationType: "invalid", Amount: 50},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Unsupported operation type: expected operation_type to be \"deposit\" or \"withdraw\"\n",
		},
		{
			name:         "Invalid amount",
			walletID:     validUUID,
			mockBalance:  100,
			mockError:    nil,
			op:           operation{OperationType: deposit, Amount: -10},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Amount must be greater than zero\n",
		},
		{
			name:         "Unexpected error",
			walletID:     validUUID,
			mockBalance:  0,
			mockError:    errors.New("unexpected error"),
			op:           operation{OperationType: deposit, Amount: 50},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to get wallet balance\n",
		},
	}

	for _, tc := range tests {
		mockDB := &mocks.DBTX{
			Balance: tc.mockBalance,
			Err:     tc.mockError,
		}

		t.Run(tc.name, func(t *testing.T) {
			// Create a new application with the mock queries
			app := &application{
				queries: database.New(mockDB),
			}

			// Create and configure a new HTTP request
			body, err := json.Marshal(tc.op)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
			req := httptest.NewRequest("POST", "/api/v1/wallets/"+tc.walletID, bytes.NewBuffer(body))
			req.SetPathValue("wallet_id", tc.walletID)
			req.Header.Set("Content-Type", "application/json")

			// Create a new response recorder
			w := httptest.NewRecorder()

			// Call the handler
			app.handleOperation(w, req)

			// Check the response
			res := w.Result()
			assert.Equal(t, tc.expectedCode, res.StatusCode)

			body, err = io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			defer res.Body.Close()

			assert.Equal(t, tc.expectedBody, string(body))
		})
	}
}
