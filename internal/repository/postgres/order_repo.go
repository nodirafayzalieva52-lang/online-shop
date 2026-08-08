package postgres

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/models"
	"shop/internal/repository"
)

// OrderRepository управляет хранением заказов в PostgreSQL
type OrderRepository struct {  
	db *sql.DB
	repo repository.OrderRepository
}

// NewOrderRepository конструктор репозитория
func NewOrderRepository(db *sql.DB, r repository.OrderRepository) *OrderRepository {
	return &OrderRepository{db: db,
	repo: r,}
}

// CreateOrder сохраняет заказ и его позиции в базе данных в единой транзакции
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	// открывает транзакцию в PostgreSQL
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Создаем запись заказа в таблице orders и получаем его ID и время создания
	orderQuery := `INSERT INTO orders (customer_id, store_id, total_price, status)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err = tx.QueryRowContext(ctx, orderQuery, order.CustomerID, order.StoreID, order.TotalPrice, order.Status).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return err
	}

	// 2. Вставляем позиции заказа в таблицу order_items
	itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price)
	VALUES ($1, $2, $3, $4) RETURNING id`

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		err = tx.QueryRowContext(ctx, itemQuery,
			order.Items[i].OrderID,
			order.Items[i].ProductID,
			order.Items[i].Quantity,
			order.Items[i].Price,
		).Scan(&order.Items[i].ID)

		if err != nil {
			return err
		}
	}

	// Подтверждаем сохранение всех данных
	return tx.Commit()
}

// GetByCustomerID возвращает список всех заказов конкретного покупателя
func (r *OrderRepository) GetByCustomerID(ctx context.Context, customerID int) ([]models.Order, error) {
	query := `SELECT id, customer_id, store_id, total_price, status, created_at
	FROM orders WHERE customer_id = $1 ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.StoreID,
			&o.TotalPrice,
			&o.Status,
			&o.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) GetByStoreID(ctx context.Context, storeID int) ([]models.Order, error) {
	query := `SELECT id, customer_id, store_id, total_price, status, created_at FROM orders WHERE store_id = $1 ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.StoreID,
			&o.TotalPrice,
			&o.Status,
			&o.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*models.Order, error) {
	query := `SELECT id, customer_id, store_id, total_price, status, created_at FROM orders WHERE id = $1`

	var o models.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID,
		&o.CustomerID,
		&o.StoreID,
		&o.TotalPrice,
		&o.Status,
		&o.CreatedAt,
	)
	
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}

	return &o, nil
}
