package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"shop/internal/delivery/http/dto"
	"shop/internal/delivery/http/middleware"
	"shop/internal/domain"
	"shop/internal/service"
	appErr "shop/pkg/errors"
	"shop/pkg/logger"
)

type OrderHandler struct {
	OrderService *service.OrderService
	log          *logger.Logger
}

func NewOrderHandler(orderService *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
		log:          log,
	}
}

// @Summary Create Order
// @Description Create Order
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body dto.CreateOrderRequest true "Order Data"
// @Success 201 {object} domain.Order
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [POST]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if len(req.Items) == 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "order must contain at least one item")
		return
	}

	orderItems := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		storeID := item.StoreID
		if storeID == 0 {
			storeID = req.StoreID
		}
		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			StoreID:   storeID,
			Quantity:  item.Quantity,
		}
	}

	order := &domain.Order{
		CustomerID: customerID,
		StoreID:    req.StoreID,
		Items:      orderItems,
	}

	if err := h.OrderService.Create(r.Context(), order); err != nil {
		if errors.Is(err, appErr.ErrInsufficientStock) {
			RespondWithError(w, http.StatusConflict, "INSUFFICIENT_STOCK", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrProductNotFound) || errors.Is(err, appErr.ErrInvalidStore) {
			RespondWithError(w, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, order)
}

// @Summary Get Order
// @Description Get order by ID
// @Tags Order
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id} [GET]
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid order id")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	order, err := h.OrderService.GetByID(r.Context(), id, userID, role)
	if err != nil {
		if errors.Is(err, appErr.ErrOrderNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrAccessDenied) {
			RespondWithError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, order)
}

// @Summary Get User Orders
// @Description Get all orders for a specific user
// @Tags Order
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [GET]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	orders, err := h.OrderService.GetByCustomerID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	if orders == nil {
		orders = []*domain.Order{}
	}

	RespondWithJSON(w, http.StatusOK, orders)
}

// @Summary Update Order Status
// @Description Update order status
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param body dto.UpdateOrderStatusRequest true "Order Status Update Data"
// @Success 200 {object} domain.Order
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id}/status [PATCH]
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	role, _ := middleware.GetRoleFromContext(r.Context())
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
		return
	}

	id, err := ParseIDParam(r, "id")
	if err != nil || id <= 0 {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid order id")
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Status != domain.OrderStatusPending && req.Status != domain.OrderStatusPaid && req.Status != domain.OrderStatusCancelled {
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid order status value")
		return
	}

	err = h.OrderService.UpdateStatus(r.Context(), id, userID, role, req.Status)
	if err != nil {
		if errors.Is(err, appErr.ErrAccessDenied) {
			RespondWithError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrOrderNotFound) {
			RespondWithError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, appErr.ErrInvalidStatusTransition) {
			RespondWithError(w, http.StatusConflict, "INVALID_TRANSITION", err.Error())
			return
		}
		RespondWithError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "order status updated successfully"})
}
