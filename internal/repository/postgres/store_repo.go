package postgres

import (
	"context"
	"database/sql"
	"errors"

	"shop/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreRepo struct {
	db *pgxpool.Pool
}

func NewStoreRepository(db *pgxpool.Pool) *StoreRepo{
	return &StoreRepo{db: db,
	}
}

func (r *StoreRepo) Create(ctx context.Context, store models.Store) error {
	query := `INSERT INTO stores (seller_id, name, description, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, store.Seller_ID, store.Name, store.Description).Scan(store.ID, store.CreatedAt)
}

func (r *StoreRepo) GetByID(ctx context.Context, id int) (*models.Store, error) {
	query := `SELECT id, seller_id, namme, description, created_at FROM stores WHERE id = $1`

	var store models.Store
	err := r.db.QueryRow(ctx, query, id).Scan(&store.ID, &store.Seller_ID, &store.Name, &store.Description, &store.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepo) GetBySellerID(ctx context.Context, sellerID int) (*models.Store, error) {
	query := `SELCT id, seller_id, name, descroption, created_at FROM stores WHERE seller_id = $1`

	var store models.Store
	err := r.db.QueryRow(ctx, query, sellerID).Scan(&store.ID, &store.Seller_ID, &store.Name, &store.Description, &store.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}
