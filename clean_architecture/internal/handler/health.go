package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Application struct {
	AppName string
}

func (a *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	appName := a.AppName
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Sistem Ayakta (GET), %v", appName)))

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
			Message: fmt.Sprintf("Basariyla Giris Yapıldı %v", appName),
			Status:  "OK-OK",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(resp)
	default:
		http.Error(w, "Meth Not Allowed", http.StatusMethodNotAllowed)
	}
}

type AuthRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthResponse struct {
	Message string
	Status  string
}
