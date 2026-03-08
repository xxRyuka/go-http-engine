package domain

import "time"

type User struct {
	ID        int
	Email     string
	Password  string // hash
	CreatedAt time.Time
}

type UserRepository interface {
	Create(u *User) error
	GetByMail(email string) (*User, error)
}

// Soru 1 : User Repodaki create ve Registerin farkı ne olacak tam olarak ? repo ve service dedik ama neden 2side burda ? amaç ne tam onu ayıkamadım
type AuthService interface {
	Register(email string, password string) error
	Login(email string, password string) (string, error) // ilerde JWT Token donecek
}
