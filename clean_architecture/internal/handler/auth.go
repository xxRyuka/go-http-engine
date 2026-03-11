package handler

import (
	"clean_architecture/internal/domain"
	"clean_architecture/internal/middleware"
	"clean_architecture/pkg/response"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

//type GetProfileResponse struct {
//
//}

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
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			// HTTP 409 Conflict dönüyoruz. Veritabanı detayı (SQLSTATE vs.) ASLA sızdırılmaz.
			response.Error(w, http.StatusConflict, "Bu e-posta adresi zaten kullanimda.")
			return // Kodun çalışmasını burada kes!
		}
		//http.Error(w, fmt.Sprintf("hata : %v", err), http.StatusNotFound)
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("register hatasi,err : %v", err))

		return
	}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	resp := LoginResponse{
		Token:   fmt.Sprintf("Dumenden-Token-%v", req.Email),
		Message: "Basariyla Kayit Olundu",
	}

	response.Json(w, http.StatusOK, resp)
	//json.NewEncoder(w).Encode(resp)

}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Json Formati Hatali")
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		// Not: Burada ham hatayı dönüyoruz, ancak production'da bu hatalar
		// formatlanıp (translator ile) kullanıcı dostu JSON'a çevrilmelidir.
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("err : %v", err.Error()))
		return
	}
	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Error(w, http.StatusBadRequest, fmt.Sprintf("Gecersiz Kullanici Adi Sifre"))
			return

		}
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("err : %v", err.Error()))

		return
	}

	resp := LoginResponse{
		Token:   token,
		Message: "Basariyla Giris Yapıldı",
	}

	response.Json(w, http.StatusOK, resp)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(middleware.UserIDKey)
	stringVal, ok := val.(string)
	if !ok {
		response.Error(w, 500, "Sunucu hatası: Context icinde UserID bulunamadi")
		return
	}
	id, err := strconv.Atoi(stringVal)
	if err != nil {
		response.Error(w, 400, "Gecersiz id formati")

		return
	}

	response.Json(w, 200, id)
}
