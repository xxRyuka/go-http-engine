package middleware

import (
	"clean_architecture/pkg/response"
	"clean_architecture/pkg/security"
	"context"
	"errors"
	"net/http"
	"strings"
)

// Mimari Standart: Context anahtarlarının çakışmaması için
// dışarıya kapalı (unexported) özel bir tip tanımlıyoruz.
type authContextKey string

// Dışarıdan Controller'ların (Handler) bu veriyi okuyabilmesi için
// anahtarı sabit (const) olarak dışa açıyoruz (exported).
const UserIDKey authContextKey = "userID"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(writer, 401, "authHeaderin ici bos")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(writer, 401, "authHeaderin formati yanlis dostum 'Bearer <token>' bekleniliyor ")
			return
		}

		tokenString := parts[1]

		claims, err := security.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, security.ErrExpiredToken) {
				response.Error(writer, 401, "tokenin süresi dolmus")
				return
			}

			response.Error(writer, 401, "gecersiz token")
			return

		}

		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), UserIDKey, claims.Subject)))
	})
}
