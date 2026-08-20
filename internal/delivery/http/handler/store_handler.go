package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"shop/internal/delivery/http/dto"
	"shop/internal/delivery/http/middleware"
	"shop/internal/service"
	appErr "shop/pkg/errors"
)

type StoreHandler struct {
	storeService *service.StoreService
}

func NewStoreHandler(storeService *service.StoreService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

// CreateStore
// 
// @Summary Create Store
// @Description Create Store
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Product json
// @Param body dto.CreateStoreRequest true "Store Data"
// @Success 201 {object} domain.Store
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /stores [POST]
func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	var req dto.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	store, err := h.storeService.CreateStore(r.Context(), userID, req.Name, req.Description)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, store)
}

// GetStore
// @Summery Get Store
// @Description Get store by ID
// @Tags Store
// @security BearerAuth
// @Produce json
// @Param id path int true "Store ID"
// @Success 200 {object} OK
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /stores/{id} [GET]
func (h *StoreHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid store id")
		return
	}

	store, err := h.storeService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErr.ErrStoreNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, store)
}

// GetBySellerID
// 
// @Summary Get store by seller ID
// @Description Get store details by seller ID path parameter or from authenticated user context
// @Tags Store
// @Security BearerAuth
// @Produce json
// @Param seller_id path int false "Seller ID (optional if authenticated)"
// @Success 200 {object} domain.Store
// @Failure 400 {object} Bad Request
// @Failure 404 {object} Not Found
// @Failure 500 {object} Internal Server Error 
// @Router /stores/seller [GET]
func (h *StoreHandler) GetBySellerID(w http.ResponseWriter, r *http.Request) {
	sellerID, err := ParseIDParam(r, "seller_id")
	if err != nil || sellerID <= 0 {
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		if ok {
			sellerID = userID
		} else {
			RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid seller id")
			return
		}
	}

	store, err := h.storeService.GetBySellerID(r.Context(), sellerID)
	if err != nil {
		if errors.Is(err, appErr.ErrStoreNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, store)
}

// UpdateStore
// 
// @Summary Update Store
// @Description Update Store details
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body dto.UpdateStoreRequest true "Store Update Data"
// @Success 200 {object} domain.Store
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /stores/{id} [PATCH]
func (h *StoreHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid store id")
		return
	}

	var req dto.UpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	store, err := h.storeService.UpdateStore(r.Context(), id, userID, role, req.Name, req.Description)
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

	RespondWithJSON(w, http.StatusOK, store)
}

// DeleteStore
// 
// @Summary Delete Store
// @Description Delete store by ID
// @Tags Store
// @Security BearerAuth
// @Produce json
// @Param id path int true "Store ID"
// @Success 200 {object} OK
// @Failure 400 {object} Bad Request
// @Failure 500 {object} Internal Server Error
// @Router /stores/{id} [DELETE]
func (h *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid store id")
		return
	}

	err = h.storeService.DeleteStore(r.Context(), id, userID, role)
	if err != nil {
		if errors.Is(err, appErr.ErrAccessDenied) {
			RespondWithError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrStoreNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "store deleted successfully"})
}
