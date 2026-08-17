package service

import (
	"context"
	"fmt"
	

	"shop/internal/models"
	"shop/internal/repository"
	"shop/pkg/jwt"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(UserRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: UserRepo,
	}
}

func (s *UserService) Register(ctx context.Context, email, password string) error {
	user := &models.User{
		Email:    email,
		PasswordHash: password,
		Role: "customer",
	}

	err := s.UserRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("s.UserRepo.GetByEmail: %w", err)
	}
	if user.PasswordHash != password {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.GenerateToken(int(user.ID), user.Email, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("token creation failed: %w", err)
	}

	return token, nil
}