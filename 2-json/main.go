package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	Message string
	Status  string
}

func handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("System Up (GET)"))

	case http.MethodPost:

		var req AuthRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			msg := fmt.Sprintf("Hatali Alındı Jsonu Decode Edemedi (POST), err : %v", err)
			//w.WriteHeader(http.StatusUnauthorized)
			//w.Write([]byte(msg))
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		resp := AuthResponse{
			Message: "Basariyla Giris Yapıldı",
			Status:  "OK-OK",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(resp)
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
