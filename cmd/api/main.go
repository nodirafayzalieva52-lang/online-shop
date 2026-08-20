package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shop/internal/config"
	"shop/internal/delivery/http/handler"
	"shop/internal/delivery/http/middleware"
	"shop/internal/delivery/http/router"
	"shop/internal/repository/postgres"
	"shop/internal/service"
	"shop/pkg/jwt"
	"shop/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// @title           Online Shop API
// @version         1.0
// @description     API сервис интернет-магазина
// @host            localhost:8080
// @BasePath        /cmd/api
func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.env"
	}

	cfg, err := config.New(configPath)
	if err != nil {
		log.Printf("warning: cleanenv.ReadConfig failed (%v), falling back to env variables", err)
		cfg = &config.Config{}
	}

	appLogger, err := logger.New(true)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer appLogger.Sync()

	ttl, err := time.ParseDuration(cfg.JWT.TTL)
	if err != nil {
		ttl = 24 * time.Hour
	}

	jwtService, err := jwt.NewService(cfg.JWT.Secret, ttl)
	if err != nil {
		appLogger.Fatal("failed to init jwt service", zap.Error(err))
	}

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User,
		cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode,
	)

	pool, err := pgxpool.New(dbCtx, connString)
	if err != nil {
		appLogger.Fatal("failed to connect to database pool", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(dbCtx); err != nil {
		appLogger.Warn("database ping failed on startup", zap.Error(err))
	} else {
		appLogger.Info("database connection established successfully")
	}

	// ---------- REPOSITORIES ----------
	userRepo := postgres.NewUserRepository(pool)
	storeRepo := postgres.NewStoreRepository(pool)
	categoryRepo := postgres.NewCategoryRepository(pool)
	productRepo := postgres.NewProductRepository(pool)
	orderRepo := postgres.NewOrderRepository(pool)

	// ---------- SERVICES ----------
	userService := service.NewUserService(userRepo, jwtService)
	storeService := service.NewStoreService(storeRepo, userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, storeRepo)
	orderService := service.NewOrderService(orderRepo, storeRepo)

	// ---------- HANDLERS ----------
	userHandler := handler.NewUserHandler(userService, appLogger)
	storeHandler := handler.NewStoreHandler(storeService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService, appLogger)
	healthHandler := handler.NewHealthHandler(pool)

	// ---------- ROUTER ----------
	r := router.NewRouter(router.Deps{
		UserHandler:     userHandler,
		StoreHandler:    storeHandler,
		CategoryHandler: categoryHandler,
		ProductHandler:  productHandler,
		OrderHandler:    orderHandler,
		HealthHandler:   healthHandler,
		AuthMiddleware:  authMiddleware,
		Logger:          appLogger,
		AllowedOrigins:  "*",
	})

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ---------- GRACEFUL SHUTDOWN ----------
	shutdownErr := make(chan error, 1)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		appLogger.Info("shutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownErr <- server.Shutdown(shutdownCtx)
	}()

	appLogger.Info("server is running", zap.String("addr", cfg.HTTPPort))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Fatal("server failed to listen and serve", zap.Error(err))
	}

	if err := <-shutdownErr; err != nil {
		appLogger.Error("server shutdown returned error", zap.Error(err))
	} else {
		appLogger.Info("server stopped gracefully")
	}
}
