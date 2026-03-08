package handler

import (
	"clean_architecture/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}
type AuthHandler struct {
	svc domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{svc: authService}
}

// Burası Aslında Controllerimiz Bizim
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Meth Not Allowed", http.StatusMethodNotAllowed)
		return

	}
	var req RegisterRequest
	s := json.NewDecoder(r.Body)
	err := s.Decode(&req)
	if err != nil {
		http.Error(w, "Json Formatı Hatalı", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Düzgün Veri Gonder Yorma Beni", http.StatusBadRequest)
		return

	}

	err = h.svc.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("hata : %v", err), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := LoginResponse{
		Token:   fmt.Sprintf("Dumenden-Token-%v", req.Email),
		Message: "Basariyla Kayit Olundu",
	}
	json.NewEncoder(w).Encode(resp)

}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Meth Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Json Formatı Hatalı", http.StatusBadRequest)
		return
	}

	//Bunu 2 kere Yazdik bu büyük bir sorun her zaman verilerin doluluğunu boslugunu boyle kontrol edemeyiz ve bazen bazıl alanların zorunlu olması bazılarının opsiyonel olması gerekiyor
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Düzgün Veri Gonder Yorma Beni", http.StatusBadRequest)

	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("hata : %v", err), http.StatusNotFound)

		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	resp := LoginResponse{
		Token:   token,
		Message: "Basariyla Giris Yapıldı",
	}

	json.NewEncoder(w).Encode(resp)
}
