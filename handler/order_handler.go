package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"shop/internal/models"
	"shop/internal/service"
	"shop/pkg/logger"
)

type OrderHandler struct {
	OrderService *service.OrderService
	log          *logger.Logger
}

func NewOrderHandler(OrderService *service.OrderService, log *log.Logger) *OrderHandler {
	return &OrderHandler{
		OrderService: OrderService,
		log: &logger.Logger{},
	}
}

type CreateOrderRequest struct {
	UserID     int     `json:"user_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный формат JSON", http.StatusBadRequest)
		return
	}

	order := h.OrderService.CreateOrder(r.Context(), &models.Order{})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Некорректный ID заказа", http.StatusBadRequest)
		return
	}

	order, err := h.OrderService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}