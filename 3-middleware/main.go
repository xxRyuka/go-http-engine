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

// handler alip handler donuyoz
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {

		start := time.Now()

		// O işini yapıp writter içine içine JSON basana kadar burası bekler.
		next(writer, request)

		total := time.Since(start)

		// Asıl handler işini bitirdi. Artık bitiş süresini hesaplayabiliriz.
		fmt.Printf("Gelen İstek: %s %s, %v Süre İçinde İslendi \n", request.Method, request.URL.Path, total.Microseconds())

	}
}
func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", LoggingMiddleware(handler))

	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
