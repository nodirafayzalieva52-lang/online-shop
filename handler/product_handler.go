package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shop/internal/service"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(ProductService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: ProductService,
	}
}

type CreateProductRequest struct {
	StoreID     int     `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Incorect JSON", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.CreateProduct(
		r.Context(),
		req.StoreID,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Incorect ID of product", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetByStoreID(w http.ResponseWriter, r *http.Request) {
	storeIDStr := r.URL.Query().Get("store_id")
	storeID, err := strconv.Atoi(storeIDStr)
	if err != nil || storeID <= 0 {
		http.Error(w, "Incorect ID of magazin", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.GetByStoreID(r.Context(), storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
