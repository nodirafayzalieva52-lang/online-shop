package postgres

import (
	"context"
	"errors"

	"shop/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StoreRepo implements repository.StoreRepository for PostgreSQL.
type StoreRepo struct {
	db *pgxpool.Pool
}

// NewStoreRepository constructs a new StoreRepo instance.
func NewStoreRepository(db *pgxpool.Pool) *StoreRepo {
	return &StoreRepo{db: db}
}

// Create inserts a new store record into database.
func (r *StoreRepo) Create(ctx context.Context, store *domain.Store) error {
	query := `INSERT INTO stores (seller_id, name, description)
	VALUES ($1, $2, $3) RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, store.SellerID, store.Name, store.Description).Scan(&store.ID, &store.CreatedAt)
}

// GetByID fetches a store by primary key.
func (r *StoreRepo) GetByID(ctx context.Context, id int64) (*domain.Store, error) {
	query := `SELECT id, seller_id, name, description, created_at FROM stores WHERE id = $1`

	var store domain.Store
	err := r.db.QueryRow(ctx, query, id).Scan(&store.ID, &store.SellerID, &store.Name, &store.Description, &store.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

// GetBySellerID fetches a store by seller's user ID.
func (r *StoreRepo) GetBySellerID(ctx context.Context, sellerID int64) (*domain.Store, error) {
	query := `SELECT id, seller_id, name, description, created_at FROM stores WHERE seller_id = $1`

	var store domain.Store
	err := r.db.QueryRow(ctx, query, sellerID).Scan(&store.ID, &store.SellerID, &store.Name, &store.Description, &store.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

// Update modifies store name and description.
func (r *StoreRepo) Update(ctx context.Context, store *domain.Store) error {
	query := `UPDATE stores SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, store.Name, store.Description, store.ID)
	return err
}

// Delete removes a store by ID.
func (r *StoreRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM stores WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
