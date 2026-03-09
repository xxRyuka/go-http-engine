package repository

import (
	"clean_architecture/internal/domain"
	"database/sql"
)

type postgresUserRepository struct {
	db *sql.DB // Bağlantı havuzunu (Connection Pool) içinde tutar
}

func NewPostgresUserRepository(db *sql.DB) *postgresUserRepository {
	return &postgresUserRepository{db: db}
}

// neden parametre olan userimiz bir pointer ?
func (r *postgresUserRepository) Create(user *domain.User) error {
	r.db.QueryRow("insert into users(email,password_hash) values ($1,$2) returning id,created_at",
		user.Email, user.Password).Scan(&user.ID, &user.CreatedAt) // burda neden & kullandık

	//scan tam olarak ne ise yariyor ?
	return nil
}

func (r *postgresUserRepository) GetByMail(email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
