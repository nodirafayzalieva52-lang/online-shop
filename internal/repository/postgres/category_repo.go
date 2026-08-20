package postgres

import (
	"context"
	"errors"

	"shop/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryRepo implements repository.CategoryRepository for PostgreSQL.
type CategoryRepo struct {
	db *pgxpool.Pool
}

// NewCategoryRepository constructs a new CategoryRepo instance.
func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// Create inserts a new category record into database.
func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`

	return r.db.QueryRow(ctx, query, category.Name).Scan(&category.ID)
}

// GetAll returns all registered categories ordered by name.
func (r *CategoryRepo) GetAll(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetByID fetches a category by primary key.
func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`

	var cat domain.Category
	err := r.db.QueryRow(ctx, query, id).Scan(&cat.ID, &cat.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// Update modifies a category name.
func (r *CategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, category.Name, category.ID)
	return err
}

// Delete removes a category by ID.
func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
