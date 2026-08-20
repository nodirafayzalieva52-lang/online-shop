package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shop/internal/delivery/http/dto"
)

// RespondWithJSON writes a JSON response with status code and payload.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

// RespondWithError writes a standardized JSON error response.
func RespondWithError(w http.ResponseWriter, code int, errCode string, message string) {
	RespondWithJSON(w, code, dto.ErrorResponse{
		Error: dto.APIError{
			Code:    errCode,
			Message: message,
		},
	})
}

// ParseIDParam extracts route path param or query param as int64.
func ParseIDParam(r *http.Request, key string) (int64, error) {
	val := r.PathValue(key)
	if val == "" {
		val = r.URL.Query().Get(key)
	}
	return strconv.ParseInt(val, 10, 64)
}
