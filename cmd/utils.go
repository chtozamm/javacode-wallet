package main

import (
	"fmt"
	"log"
	"net/http"
)

func writeResponse(w http.ResponseWriter, payload any) {
	_, err := fmt.Fprintln(w, payload)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
