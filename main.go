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
	json.NewEncoder(w).Encode(res)
}

func main() {
	// 静的ファイルの提供（オプション）
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

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
