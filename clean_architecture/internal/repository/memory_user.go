package repository

import (
	"clean_architecture/internal/domain"
	"fmt"
)

type memoryUserRepository struct {
	// Verileri RAM'de tutacağımız Map.
	// Anahtar: Email (Hızlı arama için), Değer: User nesnesinin bellek adresi (*domain.User)
	users map[string]*domain.User
}

func NewMemoryUserRepository() *memoryUserRepository {
	return &memoryUserRepository{users: make(map[string]*domain.User)}
}

func (r *memoryUserRepository) Create(u *domain.User) error {

	generatedID := len(r.users) + 1
	u.ID = generatedID
	r.users[u.Email] = u
	return nil
}

func (r *memoryUserRepository) GetByMail(mail string) (*domain.User, error) {
	// Burası Logic içermemeli !!
	//if mail == "" {
	//	return nil, fmt.Errorf("Bos Mail Gonderme")
	//}

	u, exists := r.users[mail]

	if exists != true {
		return nil, fmt.Errorf("KUllanici Bulunamadi")

	}
	return u, nil
}
