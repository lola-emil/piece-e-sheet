package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo AuthRepository
}

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
}

func NewAuthService(repo AuthRepository) AuthService {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Password:    string(hashedPassword),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &AuthResponse{User: *user, Token: token}, nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &AuthResponse{User: *user, Token: token}, nil
}
