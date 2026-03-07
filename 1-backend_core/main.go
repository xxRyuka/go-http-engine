package main

import (
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("System Up (GET)"))

	case http.MethodPost:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Veri Alındı (POST)"))

	default:
		http.Error(w, "Meth Not Allowed", http.StatusMethodNotAllowed)
	}
}
func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler)

	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
