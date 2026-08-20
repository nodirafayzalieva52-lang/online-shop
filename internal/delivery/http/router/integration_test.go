package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shop/internal/delivery/http/dto"
	"shop/internal/delivery/http/handler"
	"shop/internal/delivery/http/middleware"
	"shop/internal/delivery/http/router"
	"shop/internal/domain"
	"shop/internal/service"
	"shop/pkg/jwt"
	"shop/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// In-Memory Repository Implementations for Integration Testing
type memoryUserRepo struct {
	users  map[int64]*domain.User
	emails map[string]*domain.User
	nextID int64
}

func newMemoryUserRepo() *memoryUserRepo {
	return &memoryUserRepo{
		users:  make(map[int64]*domain.User),
		emails: make(map[string]*domain.User),
		nextID: 1,
	}
}

func (r *memoryUserRepo) Create(ctx context.Context, user *domain.User) error {
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	r.emails[user.Email] = user
	return nil
}

func (r *memoryUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return r.users[id], nil
}

func (r *memoryUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.emails[email], nil
}

func (r *memoryUserRepo) Update(ctx context.Context, user *domain.User) error {
	r.users[user.ID] = user
	r.emails[user.Email] = user
	return nil
}

func (r *memoryUserRepo) UpdateRole(ctx context.Context, userID int64, role domain.Role) error {
	if u, ok := r.users[userID]; ok {
		u.Role = role
	}
	return nil
}

type memoryStoreRepo struct {
	stores map[int64]*domain.Store
	nextID int64
}

func newMemoryStoreRepo() *memoryStoreRepo {
	return &memoryStoreRepo{
		stores: make(map[int64]*domain.Store),
		nextID: 1,
	}
}

func (r *memoryStoreRepo) Create(ctx context.Context, store *domain.Store) error {
	store.ID = r.nextID
	r.nextID++
	r.stores[store.ID] = store
	return nil
}

func (r *memoryStoreRepo) GetByID(ctx context.Context, id int64) (*domain.Store, error) {
	return r.stores[id], nil
}

func (r *memoryStoreRepo) GetBySellerID(ctx context.Context, sellerID int64) (*domain.Store, error) {
	for _, s := range r.stores {
		if s.SellerID == sellerID {
			return s, nil
		}
	}
	return nil, nil
}

func (r *memoryStoreRepo) Update(ctx context.Context, store *domain.Store) error {
	r.stores[store.ID] = store
	return nil
}

func (r *memoryStoreRepo) Delete(ctx context.Context, id int64) error {
	delete(r.stores, id)
	return nil
}

type memoryCategoryRepo struct {
	categories map[int64]*domain.Category
	nextID     int64
}

func newMemoryCategoryRepo() *memoryCategoryRepo {
	return &memoryCategoryRepo{
		categories: make(map[int64]*domain.Category),
		nextID:     1,
	}
}

func (r *memoryCategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	category.ID = r.nextID
	r.nextID++
	r.categories[category.ID] = category
	return nil
}

func (r *memoryCategoryRepo) GetAll(ctx context.Context) ([]*domain.Category, error) {
	var list []*domain.Category
	for _, c := range r.categories {
		list = append(list, c)
	}
	return list, nil
}

func (r *memoryCategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	return r.categories[id], nil
}

func (r *memoryCategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	r.categories[category.ID] = category
	return nil
}

func (r *memoryCategoryRepo) Delete(ctx context.Context, id int64) error {
	delete(r.categories, id)
	return nil
}

type memoryProductRepo struct {
	products map[int64]*domain.Product
	nextID   int64
}

func newMemoryProductRepo() *memoryProductRepo {
	return &memoryProductRepo{
		products: make(map[int64]*domain.Product),
		nextID:   1,
	}
}

func (r *memoryProductRepo) Create(ctx context.Context, product *domain.Product) error {
	product.ID = r.nextID
	r.nextID++
	r.products[product.ID] = product
	return nil
}

func (r *memoryProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return r.products[id], nil
}

func (r *memoryProductRepo) GetByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]*domain.Product, error) {
	var list []*domain.Product
	for _, p := range r.products {
		if p.StoreID == storeID {
			list = append(list, p)
		}
	}
	return list, nil
}

func (r *memoryProductRepo) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var list []*domain.Product
	for _, p := range r.products {
		list = append(list, p)
	}
	return list, nil
}

func (r *memoryProductRepo) Update(ctx context.Context, product *domain.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *memoryProductRepo) Delete(ctx context.Context, id int64) error {
	delete(r.products, id)
	return nil
}

func (r *memoryProductRepo) UpdateStock(ctx context.Context, productID int64, delta int) error {
	if p, ok := r.products[productID]; ok {
		p.Stock += delta
	}
	return nil
}

type memoryOrderRepo struct {
	orders map[int64]*domain.Order
	nextID int64
}

func newMemoryOrderRepo() *memoryOrderRepo {
	return &memoryOrderRepo{
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

func (r *memoryOrderRepo) Create(ctx context.Context, order *domain.Order) error {
	order.ID = r.nextID
	r.nextID++
	order.Status = domain.OrderStatusPending
	order.TotalPrice = 150.0
	r.orders[order.ID] = order
	return nil
}

func (r *memoryOrderRepo) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	return r.orders[id], nil
}

func (r *memoryOrderRepo) GetByCustomerID(ctx context.Context, customerID int64) ([]*domain.Order, error) {
	var list []*domain.Order
	for _, o := range r.orders {
		if o.CustomerID == customerID {
			list = append(list, o)
		}
	}
	return list, nil
}

func (r *memoryOrderRepo) GetByStoreID(ctx context.Context, storeID int64) ([]*domain.Order, error) {
	var list []*domain.Order
	for _, o := range r.orders {
		if o.StoreID == storeID {
			list = append(list, o)
		}
	}
	return list, nil
}

func (r *memoryOrderRepo) UpdateStatus(ctx context.Context, orderID int64, newStatus domain.OrderStatus) error {
	if o, ok := r.orders[orderID]; ok {
		o.Status = newStatus
	}
	return nil
}

func TestHTTPIntegration_FullFlow(t *testing.T) {
	uRepo := newMemoryUserRepo()
	sRepo := newMemoryStoreRepo()
	cRepo := newMemoryCategoryRepo()
	pRepo := newMemoryProductRepo()
	oRepo := newMemoryOrderRepo()

	jwtSvc, err := jwt.NewService("integration-test-secret-at-least-32-bytes", 24*time.Hour)
	require.NoError(t, err)

	authMiddleware := middleware.NewAuthMiddleware(jwtSvc)
	log, _ := logger.New(true)

	uSvc := service.NewUserService(uRepo, jwtSvc)
	sSvc := service.NewStoreService(sRepo, uRepo)
	cSvc := service.NewCategoryService(cRepo)
	pSvc := service.NewProductService(pRepo, sRepo)
	oSvc := service.NewOrderService(oRepo, sRepo)

	uHandler := handler.NewUserHandler(uSvc, log)
	sHandler := handler.NewStoreHandler(sSvc)
	cHandler := handler.NewCategoryHandler(cSvc)
	pHandler := handler.NewProductHandler(pSvc)
	oHandler := handler.NewOrderHandler(oSvc, log)

	appRouter := router.NewRouter(router.Deps{
		UserHandler:     uHandler,
		StoreHandler:    sHandler,
		CategoryHandler: cHandler,
		ProductHandler:  pHandler,
		OrderHandler:    oHandler,
		AuthMiddleware:  authMiddleware,
		Logger:          log,
	})

	server := httptest.NewServer(appRouter)
	defer server.Close()

	client := server.Client()

	// 1. Register seller user
	regBody, _ := json.Marshal(dto.RegisterRequest{
		Email:    "seller@shop.com",
		Password: "password123",
	})
	resp, err := client.Post(server.URL+"/register", "application/json", bytes.NewBuffer(regBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2. Login seller user
	loginBody, _ := json.Marshal(dto.LoginRequest{
		Email:    "seller@shop.com",
		Password: "password123",
	})
	resp, err = client.Post(server.URL+"/login", "application/json", bytes.NewBuffer(loginBody))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp dto.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)
	sellerToken := loginResp.Token
	assert.NotEmpty(t, sellerToken)

	// 3. Create Store
	storeBody, _ := json.Marshal(dto.CreateStoreRequest{
		Name:        "Awesome Tech Store",
		Description: "Best tech gadgets",
	})
	req, _ := http.NewRequest("POST", server.URL+"/stores", bytes.NewBuffer(storeBody))
	req.Header.Set("Authorization", "Bearer "+sellerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var store domain.Store
	json.NewDecoder(resp.Body).Decode(&store)
	assert.Equal(t, int64(1), store.ID)

	// Promote user to admin for category creation test
	uRepo.UpdateRole(context.Background(), 1, domain.RoleAdmin)
	// Re-issue admin token
	adminToken, _ := jwtSvc.GenerateToken(1, "seller@shop.com", string(domain.RoleAdmin))

	// 4. Create Category
	catBody, _ := json.Marshal(dto.CreateCategoryRequest{
		Name: "Electronics",
	})
	req, _ = http.NewRequest("POST", server.URL+"/categories", bytes.NewBuffer(catBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 5. Create Product
	prodBody, _ := json.Marshal(dto.CreateProductRequest{
		StoreID:     1,
		CategoryID:  1,
		Name:        "Wireless Headphones",
		Description: "Noise cancelling headphones",
		Price:       150.0,
		Stock:       20,
	})
	req, _ = http.NewRequest("POST", server.URL+"/products", bytes.NewBuffer(prodBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 6. Get Products List (Public)
	resp, err = client.Get(server.URL + "/products")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 7. Create Order
	orderBody, _ := json.Marshal(dto.CreateOrderRequest{
		StoreID: 1,
		Items: []dto.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
		},
	})
	req, _ = http.NewRequest("POST", server.URL+"/orders", bytes.NewBuffer(orderBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 8. Get User Orders
	req, _ = http.NewRequest("GET", server.URL+"/orders", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
