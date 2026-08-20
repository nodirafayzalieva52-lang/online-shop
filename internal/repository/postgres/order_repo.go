package postgres

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/domain"
	appErr "shop/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepo implements repository.OrderRepository for PostgreSQL.
type OrderRepo struct {
	db *pgxpool.Pool
}

// NewOrderRepository constructs a new OrderRepo instance.
func NewOrderRepository(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create inserts order & order_items transactionally while validating product, store, stock, and overriding price.
func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var calculatedTotalPrice float64

	for i := range order.Items {
		item := &order.Items[i]

		if item.Quantity <= 0 {
			return fmt.Errorf("%w: item quantity must be greater than zero", appErr.ErrInvalidOrder)
		}

		// Lock product row for update to check stock, store ownership, and fetch official price
		productQuery := `SELECT store_id, price, stock FROM products WHERE id = $1 FOR UPDATE`
		var dbStoreID int64
		var dbPrice float64
		var dbStock int

		err := tx.QueryRow(ctx, productQuery, item.ProductID).Scan(&dbStoreID, &dbPrice, &dbStock)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w: product ID %d", appErr.ErrProductNotFound, item.ProductID)
			}
			return fmt.Errorf("failed to fetch product: %w", err)
		}

		// Check if product belongs to specified order store
		targetStoreID := item.StoreID
		if targetStoreID == 0 {
			targetStoreID = order.StoreID
		}
		if dbStoreID != targetStoreID {
			return fmt.Errorf("%w: product ID %d does not belong to store ID %d", appErr.ErrInvalidStore, item.ProductID, targetStoreID)
		}

		// Check stock sufficiency
		if dbStock < item.Quantity {
			return fmt.Errorf("%w: product ID %d has stock %d, requested %d", appErr.ErrInsufficientStock, item.ProductID, dbStock, item.Quantity)
		}

		// Deduct stock
		updateStockQuery := `UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`
		_, err = tx.Exec(ctx, updateStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		// Override item price with actual database product price
		item.Price = dbPrice
		item.StoreID = targetStoreID
		calculatedTotalPrice += dbPrice * float64(item.Quantity)
	}

	order.TotalPrice = calculatedTotalPrice
	order.Status = domain.OrderStatusPending

	orderQuery := `INSERT INTO orders (customer_id, store_id, total_price, status)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err = tx.QueryRow(ctx, orderQuery, order.CustomerID, order.StoreID, order.TotalPrice, order.Status).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	itemQuery := `INSERT INTO order_items (order_id, product_id, store_id, quantity, price)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		err = tx.QueryRow(ctx, itemQuery,
			order.Items[i].OrderID,
			order.Items[i].ProductID,
			order.Items[i].StoreID,
			order.Items[i].Quantity,
			order.Items[i].Price,
		).Scan(&order.Items[i].ID)

		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetByCustomerID fetches all orders placed by a customer.
func (r *OrderRepo) GetByCustomerID(ctx context.Context, customerID int64) ([]*domain.Order, error) {
	query := `SELECT id, customer_id, store_id, total_price, status, created_at
	FROM orders WHERE customer_id = $1 ORDER BY id DESC`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
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
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, o := range orders {
		items, err := r.getOrderItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
	}

	return orders, nil
}

// GetByStoreID fetches all orders associated with a store.
func (r *OrderRepo) GetByStoreID(ctx context.Context, storeID int64) ([]*domain.Order, error) {
	query := `SELECT id, customer_id, store_id, total_price, status, created_at FROM orders WHERE store_id = $1 ORDER BY id DESC`

	rows, err := r.db.Query(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
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
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, o := range orders {
		items, err := r.getOrderItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
	}

	return orders, nil
}

// GetByID fetches a specific order by ID with its items.
func (r *OrderRepo) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `SELECT id, customer_id, store_id, total_price, status, created_at FROM orders WHERE id = $1`

	var o domain.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID,
		&o.CustomerID,
		&o.StoreID,
		&o.TotalPrice,
		&o.Status,
		&o.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	items, err := r.getOrderItems(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

// UpdateStatus modifies an order's status and restores product stock if cancelled.
func (r *OrderRepo) UpdateStatus(ctx context.Context, orderID int64, newStatus domain.OrderStatus) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var currentStatus domain.OrderStatus
	err = tx.QueryRow(ctx, `SELECT status FROM orders WHERE id = $1 FOR UPDATE`, orderID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appErr.ErrOrderNotFound
		}
		return err
	}

	// Validate status transition
	if currentStatus == domain.OrderStatusCancelled {
		return fmt.Errorf("%w: cannot change status of cancelled order", appErr.ErrInvalidStatusTransition)
	}
	if currentStatus == domain.OrderStatusPaid && newStatus == domain.OrderStatusPending {
		return fmt.Errorf("%w: cannot revert paid order to pending", appErr.ErrInvalidStatusTransition)
	}

	// Restock products if order is being cancelled
	if newStatus == domain.OrderStatusCancelled && currentStatus != domain.OrderStatusCancelled {
		items, err := r.getOrderItems(ctx, orderID)
		if err != nil {
			return err
		}
		for _, item := range items {
			_, err = tx.Exec(ctx, `UPDATE products SET stock = stock + $1, updated_at = NOW() WHERE id = $2`, item.Quantity, item.ProductID)
			if err != nil {
				return err
			}
		}
	}

	_, err = tx.Exec(ctx, `UPDATE orders SET status = $1 WHERE id = $2`, newStatus, orderID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrderRepo) getOrderItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	query := `SELECT id, order_id, product_id, store_id, quantity, price FROM order_items WHERE order_id = $1`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.StoreID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
