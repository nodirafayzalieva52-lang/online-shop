package service_test

import (
	"context"
	"testing"
	"time"

	"shop/internal/domain"
	"shop/internal/service"
	appErr "shop/pkg/errors"
	"shop/pkg/hash"
	"shop/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		user.ID = 1
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		return nil
	}
	return args.Error(0)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if u := args.Get(0); u != nil {
		return u.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if u := args.Get(0); u != nil {
		return u.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateRole(ctx context.Context, userID int64, role domain.Role) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	jwtSvc, _ := jwt.NewService("test-secret-at-least-32-bytes-long-key", time.Hour)
	userService := service.NewUserService(mockRepo, jwtSvc)

	ctx := context.Background()
	existingUser := &domain.User{ID: 1, Email: "existing@example.com"}

	mockRepo.On("GetByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	user, err := userService.Register(ctx, "existing@example.com", "password123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, appErr.ErrUserAlreadyExists)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	jwtSvc, _ := jwt.NewService("test-secret-at-least-32-bytes-long-key", time.Hour)
	userService := service.NewUserService(mockRepo, jwtSvc)

	ctx := context.Background()
	hashedPassword, _ := hash.Generate("correctpassword")
	user := &domain.User{ID: 1, Email: "user@example.com", PasswordHash: hashedPassword, Role: domain.RoleCustomer}

	mockRepo.On("GetByEmail", ctx, "user@example.com").Return(user, nil)

	token, err := userService.Login(ctx, "user@example.com", "wrongpassword")

	assert.Empty(t, token)
	assert.ErrorIs(t, err, appErr.ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterAndLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	jwtSvc, _ := jwt.NewService("test-secret-at-least-32-bytes-long-key", time.Hour)
	userService := service.NewUserService(mockRepo, jwtSvc)

	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "newuser@example.com").Return((*domain.User)(nil), nil).Once()
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := userService.Register(ctx, "newuser@example.com", "securepass123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser@example.com", user.Email)

	// Return user object containing valid password hash when Login calls GetByEmail
	mockRepo.On("GetByEmail", ctx, "newuser@example.com").Return(user, nil).Once()
	token, err := userService.Login(ctx, "newuser@example.com", "securepass123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
