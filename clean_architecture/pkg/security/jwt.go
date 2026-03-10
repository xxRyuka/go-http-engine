package security

import (
	"clean_architecture/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-gizli-gelistirme-anahtari-123")

type JwtCustomClaims struct {
	UserID               int    `json:"user_id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // Kütüphanenin standart alanlarını içeri gömüyoruz
}

func GenerateToken(user *domain.User) (string, error) {
	claims := JwtCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "clean-arch",
			Subject:   "",  // Burda user id olması gerekmiyor mu ?
			Audience:  nil, // Buraya ne yazcam ?
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Second)),
			NotBefore: nil,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // aldiği 3. parametre olan opts nedir ?
	tokenString, err := token.SignedString(jwtSecret)          // signin string metoduyla farkı ne ?
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) {
	//TODO:Validate Token methodunda kaldım yarın devam edeceğim simdilik Generate kısmını not alıp bitiriyorum
	//jwt.Parse(tokenString)
}
