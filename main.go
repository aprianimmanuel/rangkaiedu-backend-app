package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Rangkai Edu Backend Server")
	
	// Simple HTTP server for now since we're having issues with Gin
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Welcome to Rangkai Edu Backend API"}`)
	})
	
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}