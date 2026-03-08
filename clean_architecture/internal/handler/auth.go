package handler

import (
	"clean_architecture/internal/domain"
	"clean_architecture/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"

	validator2 "github.com/go-playground/validator"
)

type RegisterRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"min=6,required"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email" `
	Password string `json:"password,omitempty" validate:"min=6,required"`
}

type LoginResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}
type AuthHandler struct {
	svc       domain.AuthService
	validator *validator2.Validate
}

func NewAuthHandler(authService domain.AuthService, v *validator2.Validate) *AuthHandler {
	return &AuthHandler{svc: authService, validator: v}
}

// Burası Aslında Controllerimiz Bizim
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	//if r.Method != http.MethodPost {
	//
	//	response.Error(w, 405, "Gecersiz Method")
	//	//http.Error(w, "Meth Not Allowed", http.StatusMethodNotAllowed)
	//	return
	//
	//}
	var req RegisterRequest
	s := json.NewDecoder(r.Body)
	err := s.Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Json Formati Hatali")
		//http.Error(w, "Json Formatı Hatalı", http.StatusBadRequest)
		return
	}

	// Burayı validator yapıyor

	//if req.Email == "" || req.Password == "" {
	//	http.Error(w, "Düzgün Veri Gonder Yorma Beni", http.StatusBadRequest)
	//	return
	//
	//}
	err = h.validator.Struct(req)
	if err != nil {
		// Not: Burada ham hatayı dönüyoruz, ancak production'da bu hatalar
		// formatlanıp (translator ile) kullanıcı dostu JSON'a çevrilmelidir.
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("validasyon hatası: %v", err))
		return
	}

	err = h.svc.Register(req.Email, req.Password)
	if err != nil {
		//http.Error(w, fmt.Sprintf("hata : %v", err), http.StatusNotFound)
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("register hatasi,err : %v", err))

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := LoginResponse{
		Token:   fmt.Sprintf("Dumenden-Token-%v", req.Email),
		Message: "Basariyla Kayit Olundu",
	}

	response.Json(w, http.StatusOK, resp)
	//json.NewEncoder(w).Encode(resp)

}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	response.Error(w, 405, "Gecersiz Method")
	//
	//	//http.Error(w, "Meth Not Allowed", http.StatusMethodNotAllowed)
	//	return
	//}
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Json Formati Hatali")
		//http.Error(w, "Json Formatı Hatalı", http.StatusBadRequest)
		return
	}

	////Bunu 2 kere Yazdik bu büyük bir sorun her zaman verilerin doluluğunu boslugunu boyle kontrol edemeyiz ve bazen bazıl alanların zorunlu olması bazılarının opsiyonel olması gerekiyor
	//if req.Email == "" || req.Password == "" {
	//	http.Error(w, "Düzgün Veri Gonder Yorma Beni", http.StatusBadRequest)
	//
	//}

	err = h.validator.Struct(req)
	if err != nil {
		// Not: Burada ham hatayı dönüyoruz, ancak production'da bu hatalar
		// formatlanıp (translator ile) kullanıcı dostu JSON'a çevrilmelidir.
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("err : %v", err.Error()))
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("err : %v", err.Error()))

		//http.Error(w, fmt.Sprintf("hata : %v", err), http.StatusNotFound)

		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	resp := LoginResponse{
		Token:   token,
		Message: "Basariyla Giris Yapıldı",
	}

	response.Json(w, http.StatusOK, resp)
	//json.NewEncoder(w).Encode(resp)
}
