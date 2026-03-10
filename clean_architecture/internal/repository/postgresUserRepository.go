package repository

import (
	"clean_architecture/internal/domain"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type postgresUserRepository struct {
	db *sql.DB // Bağlantı havuzunu (Connection Pool) içinde tutar
}

func NewPostgresUserRepository(db *sql.DB) *postgresUserRepository {
	return &postgresUserRepository{db: db}
}

// neden parametre olan userimiz bir pointer ?
func (r *postgresUserRepository) Create(user *domain.User) error {
	err := r.db.QueryRow("insert into users(email,password_hash) values ($1,$2) returning id,created_at",
		user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("Beklenmeyen Hata Oluştu (id: %v) : %w ", &user.ID, err)
	} // burda neden & kullandık

	//scan tam olarak ne ise yariyor ?
	return nil
}

func (r *postgresUserRepository) GetByMail(email string) (*domain.User, error) {
	user := domain.User{}
	err := r.db.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	return &user, nil
}
