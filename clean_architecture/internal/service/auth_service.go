package service

import (
	"clean_architecture/internal/domain"
	"errors"
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

	//userEx, err := as.repo.GetByMail(email)
	//if userEx != nil {
	//
	//	return fmt.Errorf("Bu E-Posta Adresine Sahip {%v} Kullanıcı Mevcut, err : %v ", userEx.ID, err)
	//}

	user := domain.User{
		//ID:        len(as.repo.), // idsi nasıl olacak :D => Repo katmanı ayarliyo
		Email:     email,
		Password:  fmt.Sprintf("%v-hashed", password),
		CreatedAt: time.Now(),
	}

	err := as.repo.Create(&user)
	if err != nil {
		return err
	}

	return nil

}

func (as *authService) Login(email string, password string) (string, error) { // ilerde JWT Token donecek
	user, err := as.repo.GetByMail(email)

	// Bu yaklasım yanlıs mı ?
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if user.Password != fmt.Sprintf("%v-hashed", password) {
		return "", domain.ErrInvalidCredentials
	}
	return "dümenden-Jwt-Simdilik-" + user.Email, nil
}
