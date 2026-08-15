package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shop/internal/service"
)

type StoreHandler struct {
	storeService *service.StoreService
}

func NewStoreHandler(storeService *service.StoreService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

type createStoreRequest struct {
	SellerID    int    `json:"seller_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var req createStoreRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный формат JSON", http.StatusBadRequest)
		return
	}

	store, err := h.storeService.CreateStore(r.Context(), req.SellerID, req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(store)
}

func (h *StoreHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный ID магазина", http.StatusBadRequest)
		return
	}

	store, err := h.storeService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}

func (h *StoreHandler) GetBySellerID(w http.ResponseWriter, r *http.Request) {
	sellerIDStr := r.URL.Query().Get("seller_id")
	sellerID, err := strconv.Atoi(sellerIDStr)
	if err != nil || sellerID <= 0 {
		http.Error(w, "Некорректный ID продавца", http.StatusBadRequest)
		return
	}

	store, err := h.storeService.GetBySellerID(r.Context(), sellerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store)
}