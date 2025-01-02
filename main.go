package main

import (
	"encoding/json"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type FormatRequest struct {
	Code string `json:"code"`
}

type FormatResponse struct {
	FormattedCode string `json:"formattedCode"`
	Error         string `json:"error,omitempty"`
}

func setupCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(FormatResponse{Error: message})
}

func formatHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FormatRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendJSONError(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// go/format パッケージを使用してフォーマット
	formatted, err := format.Source([]byte(req.Code))
	if err != nil {
		sendJSONError(w, "Failed to format code: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := FormatResponse{FormattedCode: string(formatted)}
	json.NewEncoder(w).Encode(res)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/format", formatHandler)
	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
