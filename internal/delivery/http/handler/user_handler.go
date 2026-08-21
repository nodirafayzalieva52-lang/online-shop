package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"shop/internal/delivery/http/dto"
	"shop/internal/delivery/http/middleware"
	"shop/internal/service"
	appErr "shop/pkg/errors"
	"shop/pkg/logger"

	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
	logger      *logger.Logger
}

func NewUserHandler(userService *service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      log,
	}
}

// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User Registration Info"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /register [POST]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode register body", zap.Error(err))
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, appErr.ErrUserAlreadyExists) {
			RespondWithError(w, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	h.logger.Info("user registered successfully", zap.String("email", req.Email))
	RespondWithJSON(w, http.StatusCreated, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// @Summary Authenticate user
// @Description Authenticate user with credentials and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User Login Credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /login [POST]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode login body", zap.Error(err))
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("failed login attempt",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid email or password")
		return
	}

	h.logger.Info("user logged in successfully", zap.String("email", req.Email))
	RespondWithJSON(w, http.StatusOK, dto.LoginResponse{Token: token})
}

// @Summary Get current user profile
// @Description Get profile details of the currently authenticated user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/me [GET]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// @Summary Update current user profile
// @Description Update email or password of the currently authenticated user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "User Update Details"
// @Success 200 {object} domain.User
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /user/me [PATCH]
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	user, err := h.userService.UpdateMe(r.Context(), userID, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, appErr.ErrUserAlreadyExists) {
			RespondWithError(w, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
