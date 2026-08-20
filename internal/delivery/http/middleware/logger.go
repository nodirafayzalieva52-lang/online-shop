package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"shop/pkg/logger"

	"go.uber.org/zap"
)

type reqIDKey string

const RequestIDKey reqIDKey = "requestID"

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func generateRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(b)
}

// RequestLogger middleware logs details of each incoming HTTP request.
func RequestLogger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = generateRequestID()
			}

			w.Header().Set("X-Request-ID", reqID)
			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
			r = r.WithContext(ctx)

			sw := &statusResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(sw, r)

			duration := time.Since(start)

			log.Info("http request",
				zap.String("request_id", reqID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sw.statusCode),
				zap.Duration("duration", duration),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// GetRequestIDFromContext extracts X-Request-ID from request context.
func GetRequestIDFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
