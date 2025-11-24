package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/codepnw/mini-ecommerce/pkg/database"
)

//go:generate mockgen -source=user_repository.go -destination=mock_user_repository.go -package=userrepository

type UserRepository interface {
	// For Transactions
	Insert(ctx context.Context, db database.DBExec, input *user.User) (*user.User, error)
	SaveRefreshToken(ctx context.Context, db database.DBExec, input *user.Auth) error
	RevokedRefreshToken(ctx context.Context, db database.DBExec, token string) error

	// Use *sql.DB
	FindByID(ctx context.Context, id int64) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	ValidateRefreshToken(ctx context.Context, token string) (int64, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Insert(ctx context.Context, db database.DBExec, input *user.User) (*user.User, error) {
	m := r.domainToModel(input)
	query := `
		INSERT INTO users (email, password, first_name, last_name)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := db.QueryRowContext(
		ctx,
		query,
		m.Email,
		m.Password,
		m.FirstName,
		m.LastName,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, errs.ErrEmailAlreadyExists
		}
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	u := new(user.User)
	query := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	u := new(user.User)
	query := `
		SELECT id, email, password, role FROM users
		WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, db database.DBExec, input *user.Auth) error {
	query := `
		INSERT INTO auth (user_id, token, expires_at) VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, revoked = FALSE
	`
	_, err := db.ExecContext(ctx, query, input.UserID, input.RefreshToken, input.ExpiresAt)
	return err
}

func (r *userRepository) RevokedRefreshToken(ctx context.Context, db database.DBExec, token string) error {
	query := `UPDATE auth SET revoked = TRUE WHERE token = $1`
	res, err := db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrTokenNotFound
	}
	return nil
}

func (r *userRepository) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	var (
		userID    int64
		revoked   bool
		expiresAt time.Time
	)
	query := `SELECT user_id, revoked, expires_at FROM auth WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&userID,
		&revoked,
		&expiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errs.ErrTokenNotFound
		}
		return 0, err
	}

	if revoked {
		return 0, errs.ErrTokenRevoked
	}
	if time.Now().After(expiresAt) {
		return 0, errs.ErrTokenExpires
	}
	return userID, nil
}
