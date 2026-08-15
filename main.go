package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shop/handler"
	
	"shop/internal/repository/postgres"
	"shop/internal/service"
	"shop/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/nodirafayzalieva52-lang/cinema/api-gateway/api/handlers/middleware"
	// "github.com/nodirafayzalieva52-lang/userservice/pkg/password"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, "postgres", "20102010", "nd"))
	if err != nil {
		log.Fatal("Failed to connect top database: ", err)
	}
	defer pool.Close()
	
	log, err := logger.New(true)

	userRepo := postgres.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, log)

	storeRepo := postgres.NewStoreRepository(pool)
	storeService := service.NewStoreService(storeRepo)
	storeHandler := handler.NewStoreHandler(storeService)

	categoryRepo := postgres.NewCategoryRepository(pool)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := postgres.NewProductRepository(pool)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderRepo := postgres.NewOrderRepository(pool)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderRepo)


	mux := http.NewServeMux()

	// ---------- USER ----------
	mux.Handle("POST /register", http.HandlerFunc(userHandler.Register))
	mux.Handle("POST /login", http.HandlerFunc(userHandler.Login))

	// ---------- STORE ----------
	mux.Handle("POST /store", http.HandlerFunc(storeHandler.CreateStore))
	mux.Handle("GET /stores", http.HandlerFunc(storeHandler.GetByID))
	mux.Handle("GET /stores/seller", http.HandlerFunc(storeHandler.GetBySellerID))

	// ---------- CATEGORY ----------
	mux.Handle("POST /category", http.HandlerFunc(categoryHandler.CreateCategory))
	mux.Handle("GET /categories", http.HandlerFunc(categoryHandler.GetAll))
	mux.Handle("GET /categories/{id}", http.HandlerFunc(categoryHandler.GetByID))

	// ---------- PRODUCT ----------
	mux.Handle("POST /product", http.HandlerFunc(productHandler.CreateProduct))
	mux.Handle("GET /product", http.HandlerFunc(productHandler.GetByID))
	mux.Handle("GET /products/store", http.HandlerFunc(productHandler.GetByStoreID))

	// ---------- ORDER ----------
	mux.Handle("POST /order", http.HandlerFunc(orderHandler.CreateOrder))
	mux.Handle("GET /orders/{id}", http.HandlerFunc(orderHandler.GetByID))
}