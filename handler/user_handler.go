package handler

import (
	"encoding/json"
	"net/http"

	"shop/internal/service"
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

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode register body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		h.logger.Warn("register validation failed: missing email or password")
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("failed to register user", 
			zap.String("email", req.Email), 
			zap.Error(err),
		)
		respondWithError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	h.logger.Info("user registered successfully", zap.String("email", req.Email))
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode login body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		h.logger.Warn("login validation failed: missing email or password")
		respondWithError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("failed login attempt", 
			zap.String("email", req.Email), 
			zap.Error(err),
		)
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	h.logger.Info("user logged in successfully", zap.String("email", req.Email))
	respondWithJSON(w, http.StatusOK, LoginResponse{Token: token})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}