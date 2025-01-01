package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type FormatRequest struct {
	Code string `json:"code"`
}

type FormatResponse struct {
	FormattedCode string `json:"formattedCode"`
}

func formatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		// CORSプリフライトリクエストへの対応
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FormatRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "リクエストの読み取りに失敗しました", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "無効なJSON形式です", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("gofmt")
	cmd.Stdin = strings.NewReader(req.Code)
	formatted, err := cmd.Output()
	if err != nil {
		http.Error(w, "コードのフォーマットに失敗しました", http.StatusInternalServerError)
		return
	}

	res := FormatResponse{FormattedCode: string(formatted)}
	w.Header().Set("Content-Type", "application/json")
	// CORSヘッダーの追加
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(res)
}

func main() {
	// APIハンドラー
	http.HandleFunc("/format", formatHandler)

	// Herokuから割り当てられたポートを取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル開発用のデフォルトポート
	}
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
