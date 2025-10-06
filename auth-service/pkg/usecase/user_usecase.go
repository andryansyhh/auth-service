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
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(username string) (*dto.UserResponse, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
	jwt      *middleware.JWTManager
}

func NewUserUsecase(userRepo repository.UserRepository, jwt *middleware.JWTManager) UserUsecase {
	return &userUsecase{userRepo: userRepo, jwt: jwt}
}

func (s *userUsecase) Register(req dto.RegisterRequest) error {
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

func (s *userUsecase) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
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

	resp := &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		// CreatedBy: user.CreatedBy,
	}

	if user.UpdatedAt != nil {
		// formatted := user.UpdatedAt.Format(time.RFC3339)
		// resp.UpdatedAt = &formatted
	}
	if user.DeletedAt != nil {
		// formatted := user.DeletedAt.Format(time.RFC3339)
		// resp.DeletedAt = &formatted
	}

	return resp, nil
}
