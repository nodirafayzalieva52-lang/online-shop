package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck checks database connectivity and returns 200 OK or 503 Service Unavailable.
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		RespondWithError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database connection error")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"database": "connected",
	})
}
