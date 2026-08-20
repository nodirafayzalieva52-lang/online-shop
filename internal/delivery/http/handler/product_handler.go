package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"shop/internal/delivery/http/dto"
	"shop/internal/delivery/http/middleware"
	"shop/internal/domain"
	"shop/internal/service"
	appErr "shop/pkg/errors"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func parsePagination(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

// CreateProduct
//
// @Summary Create Product
// @Description Create Product
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body dto.CreateProductRequest true "Product Data"
// @Success 201 {object} domain.Product
// @Failure 400 {object} 
// @Failure 500 {object}
// @Router /products [POST]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	product, err := h.ProductService.CreateProduct(r.Context(), userID, role, req)
	if err != nil {
		if errors.Is(err, appErr.ErrAccessDenied) {
			RespondWithError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrStoreNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, product)
}

// GetByID
//
// @Summary Get Movie
// @Description Get movie by ID
// @Tags Movie
// @Security BearerAuth
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} MovieResponse
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /products/{id} [GET]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid product id")
		return
	}

	product, err := h.ProductService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErr.ErrProductNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetByStoreID(w http.ResponseWriter, r *http.Request) {
	storeID, err := ParseIDParam(r, "store_id")
	if err != nil || storeID <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid store id")
		return
	}

	limit, offset := parsePagination(r)
	products, err := h.ProductService.GetByStoreID(r.Context(), storeID, limit, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	if products == nil {
		products = []*domain.Product{}
	}

	RespondWithJSON(w, http.StatusOK, products)
}


func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	products, err := h.ProductService.GetAll(r.Context(), limit, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	if products == nil {
		products = []*domain.Product{}
	}

	RespondWithJSON(w, http.StatusOK, products)
}

// UpdateProduct
//
// @Summary Update Product
// @Description Update product details
// @Tags Product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body dto.UpdateProductRequest true "Product Update Data"
// @Success 200 {object} domain.Product
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /products/{id} [PATCH]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid product id")
		return
	}

	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	product, err := h.ProductService.UpdateProduct(r.Context(), id, userID, role, req)
	if err != nil {
		if errors.Is(err, appErr.ErrAccessDenied) {
			RespondWithError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrProductNotFound) || errors.Is(err, appErr.ErrStoreNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, product)
}

// DeleteProduct
//
// @Summary Delete Product
// @Description Delete product by ID
// @Tags Product
// @Security BearerAuth
// @Produce json
// @Success 200 {object} OK
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /products/{id} [DELETE]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid product id")
		return
	}

	err = h.ProductService.DeleteProduct(r.Context(), id, userID, role)
	if err != nil {
		if errors.Is(err, appErr.ErrAccessDenied) {
			RespondWithError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrProductNotFound) || errors.Is(err, appErr.ErrStoreNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "product deleted successfully"})
}
