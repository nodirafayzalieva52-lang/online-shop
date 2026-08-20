package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"shop/internal/delivery/http/dto"
	"shop/internal/domain"
	"shop/internal/service"
	appErr "shop/pkg/errors"
)

type CategoryHandler struct {
	CategoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: categoryService,
	}
}

// CreateCategory
// 
// @Summary Create Category
// @Description Create Category
// @Tags Category
// @Security BearerAuth
// @Accept json
// @Product json
// @Param body dto.CreateCategoryRequest true "Category Data"
// @Success 201 {object} domain.Category
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /categories [POST]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	category, err := h.CategoryService.CreateCategory(r.Context(), req.Name)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, category)
}

// GetCategory
// 
// @Summery Get Category 
// @Description Get category by ID
// @Tags Category
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} domain.Category
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /categories/{id} [GET]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid category id")
		return
	}

	category, err := h.CategoryService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErr.ErrCategoryNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, category)
}

// GetAll
// 
// @Summary Get All
// @Description Get All Categories
// @Tags Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} OK
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /categories [GET]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.CategoryService.GetAll(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	if categories == nil {
		categories = []*domain.Category{}
	}

	RespondWithJSON(w, http.StatusOK, categories)
}

// UpdateCategory
// 
// @Summary Update Category
// @Description Update Category details
// @Tags Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body dto.UpdateCategoryRequest true "Category Update Data"
// @Success 200 {object} domain.Category
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /categories/{id} [PATCH]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid category id")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	category, err := h.CategoryService.UpdateCategory(r.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, appErr.ErrCategoryNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, category)
}

// DeleteCategory
// 
// @Summary Delete Category
// @Description Delete category by ID
// @Tags Category
// @Security BearerAuth
// @Produce json
// @Param id path int true "Ctaegory ID"
// @Success 200 {object} OK
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /categories/{id} [DELETE]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid category id")
		return
	}

	err = h.CategoryService.DeleteCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErr.ErrCategoryNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "category deleted successfully"})
}
