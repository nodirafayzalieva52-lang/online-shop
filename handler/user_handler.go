package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"shop/internal/models"
	"shop/internal/service"
	"shop/pkg/logger"
	"strconv"

	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *logger.Logger
	UserService service.UserService
}

func NewUserHandler(logger *logger.Logger, UserService service.UserService) *UserHandler {
	return &UserHandler{
		logger: logger,
		UserService: UserService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.logger.Error("h.UserService.Create", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.UserService.UserRepo.Create(context.TODO(), &user)
	if err != nil {
		h.logger.Error("h.UserService.Create", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("h.UserService.GetByID", zap.Error(err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.UserRepo.GetByID(context.TODO(), id)
	if err != nil {
		h.logger.Error("h.UserSErvice.GetByID", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		h.logger.Error("h.UserService.GetByID", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
