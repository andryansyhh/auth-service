package usecase

import (
	"errors"
	"time"

	"github.com/andryansyhh/auth-service/pkg/domain/dto"
	"github.com/andryansyhh/auth-service/pkg/domain/model"
	"github.com/andryansyhh/auth-service/pkg/middleware"
	"github.com/andryansyhh/auth-service/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(req dto.LoginRegisterRequest) error
	Login(req dto.LoginRegisterRequest) (*dto.AuthResponse, error)
	GetProfile(username string) (*dto.UserResponse, error)
	ListUsers() ([]dto.UserResponse, error)
	UpdateUser(id int64, req dto.UpdateUserRequest) error
	DeleteUser(id int64) error
}

type userUsecase struct {
	userRepo repository.UserRepository
	jwt      *middleware.JWTManager
}

func NewUserUsecase(userRepo repository.UserRepository, jwt *middleware.JWTManager) UserUsecase {
	return &userUsecase{userRepo: userRepo, jwt: jwt}
}

func (s *userUsecase) Register(req dto.LoginRegisterRequest) error {
	_, err := s.userRepo.GetUserByUsername(req.Username)
	if err == nil {
		return errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	createdBy := "system"
	user := &model.User{
		Username:  req.Username,
		Password:  string(hashed),
		BaseModel: model.BaseModel{CreatedBy: &createdBy},
	}
	return s.userRepo.CreateUser(user)
}

func (s *userUsecase) Login(req dto.LoginRegisterRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Message: "login success",
		Token:   token,
	}, nil
}

func (s *userUsecase) GetProfile(username string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *userUsecase) ListUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.ListUsers()
	if err != nil {
		return nil, err
	}

	var resp []dto.UserResponse
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp, nil
}

func (s *userUsecase) UpdateUser(id int64, req dto.UpdateUserRequest) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.Username = req.Username
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	} else {
		user.Password = ""
	}

	return s.userRepo.UpdateUser(user)
}

func (s *userUsecase) DeleteUser(id int64) error {
	if _, err := s.userRepo.GetUserByID(id); err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.DeleteUser(id)
}
