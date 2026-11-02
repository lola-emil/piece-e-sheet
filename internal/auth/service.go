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
	// 1. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. Create user
	user := &User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Password:    string(hashedPassword),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		// In a real app, check for duplicate email error here
		return nil, err
	}

	// 3. Generate token
	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return &AuthResponse{User: *user, Token: token}, nil
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 1. Find user
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. Generate token
	token, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &AuthResponse{User: *user, Token: token}, nil
}
