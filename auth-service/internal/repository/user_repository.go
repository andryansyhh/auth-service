package repository

import (
	"github.com/andryansyhh/auth-service/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	_, err := r.db.NamedExec(`
		INSERT INTO users (username, password, created_at, created_by)
		VALUES (:username, :password, NOW(), :created_by)
	`, user)
	return err
}

func (r *userRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username=$1 AND deleted_at IS NULL", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
