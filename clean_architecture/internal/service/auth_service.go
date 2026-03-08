package service

import (
	"clean_architecture/internal/domain"
	"fmt"
	"time"
)

// Liskov Substitution Principle (LSP) & Dependency Inversion
type authService struct {
	// DİKKAT: memoryUserRepository yazmıyoruz!
	// Domain'deki sözleşmeyi (interface) koyuyoruz.
	// Bu sayede Service, Postgres mi yoksa Memory repo mu kullanıyor ASLA bilmez!
	repo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *authService {
	return &authService{repo: userRepo}
}

func (as *authService) Register(email string, password string) error {

	mail, err := as.repo.GetByMail(email)
	if err != nil {
		return fmt.Errorf("Bu E-Posta Adresine Sahip {%v} Kullanıcı Mevcut ", mail)
	}

	user := domain.User{
		//ID:        len(as.repo.), // idsi nasıl olacak :D => Repo katmanı ayarliyo
		Email:     email,
		Password:  fmt.Sprintf("%v-hashed", password),
		CreatedAt: time.Now(),
	}

	return as.repo.Create(&user)

}

func (as *authService) Login(email string, password string) (string, error) { // ilerde JWT Token donecek
	user, err := as.repo.GetByMail(email)
	if err != nil {
		return "", err
	}

	if user.Password != fmt.Sprintf("%v-hashed", password) {
		return "", fmt.Errorf("gecersiz Sifre")
	}
	return "dümenden-Jwt-Simdilik-" + user.Email, nil
}
