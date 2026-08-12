package handler

import (
	"encoding/json"
	"net/http"

	"shop/internal/service"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func NewAuthHandler(AuthService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: *AuthService,
	}
}

type RegisterRequest struct {
	Email	 string `json:"email"`
	Password string `json:"password`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorect JSON", http.StatusBadRequest)
		return
	}

	err := h.AuthService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "user successfully registered"}`))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorect JSON", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, 
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "вы успешно вышли из системы"}`))
}
