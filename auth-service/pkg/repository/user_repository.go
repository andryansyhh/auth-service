package repository

import (
	"github.com/andryansyhh/auth-service/pkg/domain/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	ListUsers() ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int64) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	query := `INSERT INTO users (username, password, created_at, created_by)
              VALUES ($1, $2, NOW(), $3)`
	stmt, err := r.db.Preparex(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password, user.CreatedBy)
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

func (r *userRepo) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE id=$1 AND deleted_at IS NULL"
	stmt, err := r.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Get(&user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) ListUsers() ([]model.User, error) {
	var users []model.User
	query := "SELECT id, username, created_at FROM users WHERE deleted_at IS NULL ORDER BY id DESC"
	stmt, err := r.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Select(&users)
	return users, err
}

func (r *userRepo) UpdateUser(user *model.User) error {
	query := `UPDATE users SET 
                 username=$1, 
                 password=COALESCE(NULLIF($2, ''), password),
								 updated_at=NOW()
              WHERE id=$3`
	stmt, err := r.db.Preparex(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password, user.ID)
	return err
}

func (r *userRepo) DeleteUser(id int64) error {
	query := "UPDATE users SET deleted_at=NOW() WHERE id=$1"
	stmt, err := r.db.Preparex(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}
