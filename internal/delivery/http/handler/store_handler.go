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
