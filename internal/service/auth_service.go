package service

import (
	"context"
	"fmt"

	"shop/internal/models"
	"shop/internal/repository"
	"shop/pkg/jwt"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(UserRepo repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: UserRepo,
	}
}

func (a *AuthService) Register(ctx context.Context, email, password string) error {
	existsUser := a.UserRepo.GetByEmail(ctx, email)
	if existsUser != nil {
		return fmt.Errorf("user with this email is already exists")
	}

	user := models.User{
		Email: email,
		PasswordHash: password,
		Role: "customer",
	}

	err := a.UserRepo.Create(ctx, &user)
	if err != nil {
		return fmt.Errorf("a.UserRepo.Create: %w", err)
	}
	return nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user := a.UserRepo.GetByEmail(ctx, email)
	
	if user.PasswordHash != password {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.GenerateToken(int(user.ID), user.Email, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("token creation failed: %w", err)
	}
	return token, nil
}
