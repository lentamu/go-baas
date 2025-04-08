package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type apiResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	addr := flag.String("addr", ":8080", "http server address")
	flag.Parse()

	logger := log.Default()

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", withLogging(logger, handleHash))
	mux.HandleFunc("/verify", withLogging(logger, handleVerify))

	server := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("Starting server on %s", *addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Error starting server: %s", err)
	}
}

func handleHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, apiResponse{Error: "Method not allowed"})
		return
	}

	if err := r.ParseForm(); err != nil {
		sendResponse(w, http.StatusBadRequest, apiResponse{Error: "Invalid request"})
		return
	}

	raw := r.Form.Get("raw")
	costStr := r.Form.Get("cost")

	if raw == "" || costStr == "" {
		sendResponse(w, http.StatusBadRequest, apiResponse{Error: "Missing raw or cost params"})
		return
	}

	cost, err := strconv.Atoi(costStr)
	if err != nil {
		sendResponse(w, http.StatusBadRequest, apiResponse{Error: "Invalid cost"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(raw), cost)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, apiResponse{Error: "Failed to generate hash"})
		return
	}

	sendResponse(w, http.StatusOK, apiResponse{Success: true, Data: string(hash)})
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, apiResponse{Error: "Method not allowed"})
		return
	}

	if err := r.ParseForm(); err != nil {
		sendResponse(w, http.StatusBadRequest, apiResponse{Error: "Invalid request"})
		return
	}

	raw := r.Form.Get("raw")
	hash := r.Form.Get("hash")

	if raw == "" || hash == "" {
		sendResponse(w, http.StatusBadRequest, apiResponse{Error: "Missing raw or hash params"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	if err != nil {
		sendResponse(w, http.StatusBadRequest, apiResponse{Error: "Invalid password"})
		return
	}

	sendResponse(w, http.StatusOK, apiResponse{Success: true})
}

func withLogging(logger *log.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Printf("Started %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		logger.Printf("Completed %s %s from %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func sendResponse(w http.ResponseWriter, status int, resp apiResponse) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
