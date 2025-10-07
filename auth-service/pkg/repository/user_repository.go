package repository

import (
	"github.com/andryansyhh/auth-service/pkg/domain/model"
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
	query := `INSERT INTO users (username, password, created_at, created_by)
              VALUES (:username, :password, NOW(), :created_by)`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user)
	return err
}

func (r *userRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	query := "SELECT * FROM users WHERE username=$1 AND deleted_at IS NULL"

	stmt, err := r.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Get(&user, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
