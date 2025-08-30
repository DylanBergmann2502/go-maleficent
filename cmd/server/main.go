// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok","message":"Go Maleficent API is running"}`)
}

func main() {
	http.HandleFunc("/health", healthHandler)

	port := ":8000"
	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
