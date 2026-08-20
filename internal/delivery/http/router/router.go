package router

import (
	"net/http"

	"shop/internal/delivery/http/handler"
	"shop/internal/delivery/http/middleware"
	"shop/internal/domain"
	"shop/pkg/logger"
)

type Deps struct {
	UserHandler     *handler.UserHandler
	StoreHandler    *handler.StoreHandler
	CategoryHandler *handler.CategoryHandler
	ProductHandler  *handler.ProductHandler
	OrderHandler    *handler.OrderHandler
	HealthHandler   *handler.HealthHandler
	AuthMiddleware  *middleware.AuthMiddleware
	Logger          *logger.Logger
	AllowedOrigins  string
}

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()

	// ---------- HEALTH CHECK ----------
	if d.HealthHandler != nil {
		mux.HandleFunc("GET /health", d.HealthHandler.HealthCheck)
		mux.HandleFunc("GET /ping", d.HealthHandler.HealthCheck)
	}

	// ---------- PUBLIC ROUTES ----------
	mux.HandleFunc("POST /register", d.UserHandler.Register)
	mux.HandleFunc("POST /login", d.UserHandler.Login)

	mux.HandleFunc("GET /stores/{id}", d.StoreHandler.GetByID)

	mux.HandleFunc("GET /categories", d.CategoryHandler.GetAll)
	mux.HandleFunc("GET /categories/{id}", d.CategoryHandler.GetByID)

	mux.HandleFunc("GET /products", d.ProductHandler.GetAll)
	mux.HandleFunc("GET /products/{id}", d.ProductHandler.GetByID)
	mux.HandleFunc("GET /stores/{store_id}/products", d.ProductHandler.GetByStoreID)

	// ---------- PROTECTED ROUTES (Require Auth) ----------
	if d.AuthMiddleware != nil {
		auth := d.AuthMiddleware.RequireAuth
		adminOnly := d.AuthMiddleware.RequireRole(domain.RoleAdmin)
		sellerOrAdmin := d.AuthMiddleware.RequireRole(domain.RoleSeller, domain.RoleAdmin)

		// User profile
		mux.Handle("GET /me", auth(http.HandlerFunc(d.UserHandler.GetMe)))
		mux.Handle("PATCH /me", auth(http.HandlerFunc(d.UserHandler.UpdateMe)))

		// Stores
		mux.Handle("POST /stores", auth(http.HandlerFunc(d.StoreHandler.CreateStore)))
		mux.Handle("POST /store", auth(http.HandlerFunc(d.StoreHandler.CreateStore)))
		mux.Handle("GET /stores/seller", auth(http.HandlerFunc(d.StoreHandler.GetBySellerID)))
		mux.Handle("PATCH /stores/{id}", auth(http.HandlerFunc(d.StoreHandler.UpdateStore)))
		mux.Handle("DELETE /stores/{id}", auth(http.HandlerFunc(d.StoreHandler.DeleteStore)))

		// Categories (Admin restricted)
		mux.Handle("POST /categories", auth(adminOnly(http.HandlerFunc(d.CategoryHandler.CreateCategory))))
		mux.Handle("POST /category", auth(adminOnly(http.HandlerFunc(d.CategoryHandler.CreateCategory))))
		mux.Handle("PATCH /categories/{id}", auth(adminOnly(http.HandlerFunc(d.CategoryHandler.UpdateCategory))))
		mux.Handle("DELETE /categories/{id}", auth(adminOnly(http.HandlerFunc(d.CategoryHandler.DeleteCategory))))

		// Products (Seller / Admin restricted for write)
		mux.Handle("POST /products", auth(sellerOrAdmin(http.HandlerFunc(d.ProductHandler.CreateProduct))))
		mux.Handle("POST /product", auth(sellerOrAdmin(http.HandlerFunc(d.ProductHandler.CreateProduct))))
		mux.Handle("PATCH /products/{id}", auth(sellerOrAdmin(http.HandlerFunc(d.ProductHandler.UpdateProduct))))
		mux.Handle("DELETE /products/{id}", auth(sellerOrAdmin(http.HandlerFunc(d.ProductHandler.DeleteProduct))))

		// Orders
		mux.Handle("POST /orders", auth(http.HandlerFunc(d.OrderHandler.CreateOrder)))
		mux.Handle("POST /order", auth(http.HandlerFunc(d.OrderHandler.CreateOrder)))
		mux.Handle("GET /orders", auth(http.HandlerFunc(d.OrderHandler.GetUserOrders)))
		mux.Handle("GET /orders/{id}", auth(http.HandlerFunc(d.OrderHandler.GetByID)))
		mux.Handle("PATCH /orders/{id}/status", auth(sellerOrAdmin(http.HandlerFunc(d.OrderHandler.UpdateStatus))))
	}

	// Apply Middlewares: CORS -> RequestLogger -> Mux
	var handler http.Handler = mux
	if d.Logger != nil {
		handler = middleware.RequestLogger(d.Logger)(handler)
	}
	handler = middleware.CORS(d.AllowedOrigins)(handler)

	return handler
}
