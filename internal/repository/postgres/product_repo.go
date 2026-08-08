package postgres

import (
	"context"
	"database/sql"
	"errors"

	"shop/internal/models"
	"shop/internal/repository"
)

type ProductRepo struct {
	db *sql.DB
	repo repository.ProductRespository
}

func NewProductRepository(db *sql.DB, r repository.ProductRespository) *ProductRepo {
	return &ProductRepo{db: db, 
	repo: r,}
}

func (r *ProductRepo) Create(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (store_id, category_id, title, description, price, stock, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query, product.StoreID, product.CategoryID, product.Title,
		product.Description, product.Price, product.Stock,).Scan(&product.ID, &product.Created_At, &product.Updated_At)
}

func (r *ProductRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `SELECT id, store_id, category_id, title, description, price, stock, created_at, updated_at
	FROM products WHERE id = $1`

	var p models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.StoreID, &p.CategoryID, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Created_At, &p.Updated_At)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, err
}

func (r *ProductRepo) GetByStoreID(ctx context.Context, storeID int) ([]*models.Product, error) {
	query := `SELECT id, store_id, category_id, title, description, price, stock, created_at, updated_at
	FROM products WHERE store_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
			var p models.Product
		if err := rows.Scan(&p.ID, &p.StoreID, &p.CategoryID, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Created_At, &p.Updated_At); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, store_id, category_id, name, description, price, stock, created_at
	FROM prodycts ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.StoreID,
			&p.CategoryID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.Created_At,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) UpdateStock(ctx context.Context, ProductID int, delta int) error {
	query := `UPDATE products SET stock = stock + $1 WHERE id = $2 AND (stock = $1) >= 0`

	result, err := r.db.ExecContext(ctx, query, delta, ProductID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("insufficient stock or product not found")
	}
	return nil
}
