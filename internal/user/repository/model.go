package userrepository

import (
	"time"

	"github.com/codepnw/mini-ecommerce/internal/user"
)

type userModel struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *userRepository) modelToDomain(m *userModel) *user.User {
	return &user.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *userRepository) domainToModel(m *user.User) *userModel {
	return &userModel{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
