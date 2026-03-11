package security

import (
	"clean_architecture/internal/domain"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// artık env dosyamiz var
//var jwtSecret = []byte("super-gizli-gelistirme-anahtari-123")

type JwtCustomClaims struct {
	UserID               int    `json:"user_id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // Kütüphanenin standart alanlarını içeri gömüyoruz
}

// Ortak kullanılacak custom hatalar tanımlıyoruz.
// Bu sayede controller/middleware katmanında hatanın tipine göre mantık kurabiliriz.
var (
	ErrInvalidToken = errors.New("token gecersiz veya uzerinde oynanmis")
	ErrExpiredToken = errors.New("token suresi dolmus")
)

func GenerateToken(user *domain.User) (string, error) {
	claims := JwtCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "clean-arch",
			Subject:   fmt.Sprintf("%v", user.ID),
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)              // aldiği 3. parametre olan opts nedir ?
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // signin string metoduyla farkı ne ?
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func keyFunc(t *jwt.Token) (any, error) {

	if t.Method != jwt.SigningMethodHS256 {
		return nil, fmt.Errorf("Gecersiz Algoritma %v", t.Method.Alg())
	}
	return []byte(os.Getenv("JWT_SECRET")), nil
}
func ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	// Adım 4: Type Assertion ve geçerlilik onayı
	// Token başarılıysa ve validasyonlardan (exp, nbf vb.) geçtiyse, struct'ımıza cast ediyoruz.
	claims, ok := token.Claims.(*JwtCustomClaims)
	if token.Valid && ok {
		return claims, nil
	}

	// her ihtimale karsı fallback
	return nil, ErrInvalidToken

}
