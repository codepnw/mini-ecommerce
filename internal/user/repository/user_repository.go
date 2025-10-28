package userrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/mini-ecommerce/internal/errs"
	"github.com/codepnw/mini-ecommerce/internal/user"
)

type UserRepository interface {
	Insert(ctx context.Context, input *user.User) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Insert(ctx context.Context, input *user.User) (*user.User, error) {
	m := r.domainToModel(input)
	query := `
		INSERT INTO users (email, password, first_name, last_name)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		m.Email,
		m.Password,
		m.FirstName,
		m.LastName,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	u := new(user.User)
	query := `SELECT email, password FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
