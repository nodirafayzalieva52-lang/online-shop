package postgres

import (
	"context"
	"errors"
	"fmt"

	"shop/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepo implements repository.ProductRepository for PostgreSQL.
type ProductRepo struct {
	db *pgxpool.Pool
}

// NewProductRepository constructs a new ProductRepo instance.
func NewProductRepository(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create inserts a new product record into database.
func (r *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (store_id, category_id, name, description, price, stock)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		product.StoreID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

// GetByID fetches a product by primary key.
func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, store_id, category_id, name, description, price, stock, created_at, updated_at
	FROM products WHERE id = $1`

	var p domain.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.StoreID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetByStoreID fetches products belonging to a store with pagination.
func (r *ProductRepo) GetByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]*domain.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `SELECT id, store_id, category_id, name, description, price, stock, created_at, updated_at
	FROM products WHERE store_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, storeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.StoreID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// GetAll fetches all products with pagination.
func (r *ProductRepo) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	query := `SELECT id, store_id, category_id, name, description, price, stock, created_at, updated_at
	FROM products ORDER BY id DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID, &p.StoreID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// Update modifies product details.
func (r *ProductRepo) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products 
	SET category_id = $1, name = $2, description = $3, price = $4, stock = $5, updated_at = NOW() 
	WHERE id = $6`
	_, err := r.db.Exec(ctx, query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID,
	)
	return err
}

// Delete removes a product by ID.
func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdateStock atomic modification of product stock.
func (r *ProductRepo) UpdateStock(ctx context.Context, productID int64, delta int) error {
	query := `UPDATE products SET stock = stock + $1, updated_at = NOW() WHERE id = $2 AND (stock + $1) >= 0`

	result, err := r.db.Exec(ctx, query, delta, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock or product not found")
	}
	return nil
}
