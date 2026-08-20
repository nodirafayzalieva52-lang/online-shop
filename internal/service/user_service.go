package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"shop/internal/domain"
	"shop/internal/repository"
	appErr "shop/pkg/errors"
	"shop/pkg/hash"
	"shop/pkg/jwt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type UserService struct {
	UserRepo   repository.UserRepository
	JWTService *jwt.Service
}

func NewUserService(userRepo repository.UserRepository, jwtService *jwt.Service) *UserService {
	return &UserService{
		UserRepo:   userRepo,
		JWTService: jwtService,
	}
}

func (s *UserService) ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (s *UserService) ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

func (s *UserService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if err := s.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := s.ValidatePassword(password); err != nil {
		return nil, err
	}

	existing, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	if existing != nil {
		return nil, appErr.ErrUserAlreadyExists
	}

	hashedPassword, err := hash.Generate(password)
	if err != nil {
		return nil, fmt.Errorf("hash.Generate: %w", err)
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         domain.RoleCustomer,
	}

	if err := s.UserRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("userRepo.Create: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("userRepo.GetByEmail: %w", err)
	}
	if user == nil {
		return "", appErr.ErrInvalidCredentials
	}
	if !hash.Compare(user.PasswordHash, password) {
		return "", appErr.ErrInvalidCredentials
	}

	token, err := s.JWTService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", fmt.Errorf("token creation failed: %w", err)
	}

	return token, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	user, err := s.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appErr.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateMe(ctx context.Context, userID int64, newEmail *string, newPassword *string) (*domain.User, error) {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if newEmail != nil && *newEmail != "" {
		if err := s.ValidateEmail(*newEmail); err != nil {
			return nil, err
		}
		if *newEmail != user.Email {
			existing, err := s.UserRepo.GetByEmail(ctx, *newEmail)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return nil, appErr.ErrUserAlreadyExists
			}
			user.Email = *newEmail
		}
	}

	if newPassword != nil && *newPassword != "" {
		if err := s.ValidatePassword(*newPassword); err != nil {
			return nil, err
		}
		hashedPassword, err := hash.Generate(*newPassword)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	if err := s.UserRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *UserService) PromoteToSeller(ctx context.Context, userID int64) error {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Role == domain.RoleCustomer {
		return s.UserRepo.UpdateRole(ctx, userID, domain.RoleSeller)
	}
	return nil
}
