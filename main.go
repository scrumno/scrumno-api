package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type Response struct {
    Message string `json:"message"`
    Status  string `json:"status"`
    Time    string `json:"time"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{
        Message: "Scrumno API is running",
        Status:  "ok",
        Time:    time.Now().Format(time.RFC3339),
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{
        Message: "Hello from Scrumno API",
        Status:  "success",
        Time:    time.Now().Format(time.RFC3339),
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/api/v1/hello", helloHandler)

    log.Println("🚀 Scrumno API Server starting on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal("❌ Server failed to start:", err)
    }
}
