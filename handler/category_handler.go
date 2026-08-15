package handler

import (
	"encoding/json"
	"net/http"
	"shop/internal/service"
	"strconv"
)

type CategoryHandler struct {
	CategoryService *service.CategoryService
}

func NewCategoryHandler(CategoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: CategoryService,
	}
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorect JSON", http.StatusBadRequest)
		return
	}

	category, err := h.CategoryService.CreateCategory(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "INcorect ID", http.StatusBadRequest)
		return
	}

	category, err := h.CategoryService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.CategoryService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
