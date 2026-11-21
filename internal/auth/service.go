package auth

import (
	"api/internal/user"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	authRepo AuthRepository
	userRepo user.UserRepository
}

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	LoginWithGoogle(ctx context.Context, info *UserInfo) (*AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*user.User, error)
}

func NewAuthService(authRepo AuthRepository, userRepo user.UserRepository) AuthService {
	return &service{authRepo: authRepo, userRepo: userRepo}
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

	err = s.authRepo.CreateUser(ctx, user)
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
	user, err := s.authRepo.FindByEmail(ctx, req.Email)
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

	refreshToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &AuthResponse{
		User:         *user,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) LoginWithGoogle(ctx context.Context, info *UserInfo) (*AuthResponse, error) {
	user, err := s.authRepo.FindByEmail(ctx, info.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user = &User{
			Email:       info.Email,
			DisplayName: info.Name,
			// add GoogleID/Picture columns to User if you want proper linking
		}
		if err := s.authRepo.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	access, err := GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	refresh, err := GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &AuthResponse{User: *user, Token: access, RefreshToken: refresh}, nil
}

func (s *service) GetMe(ctx context.Context, userID string) (*user.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
