package main

import (
	"log"
	"net/http"
)

func (app *application) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check the database connection
	err := app.db.Ping(r.Context())
	if err != nil {
		log.Printf("Database health check failed: %v\n", err)
		http.Error(w, "Database health check failed", http.StatusInternalServerError)
		return
	}

	// Write the response
	_, err = w.Write([]byte("OK\n"))
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
