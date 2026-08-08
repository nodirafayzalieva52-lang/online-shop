package postgres

import (
	"context"
	"database/sql"
	"errors"

	"shop/internal/models"
	"shop/internal/repository"
)

type CategoryRepo struct {
	db *sql.DB
	repo repository.CategoryRepository
}

func NewCategoryRepository(db *sql.DB, r repository.CategoryRepository) *CategoryRepo {
	return &CategoryRepo{db: db,
	repo: r,}
}

func (r *CategoryRepo) Create(ctx context.Context, category models.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`

	return r.db.QueryRowContext(ctx, query, category.Name).Scan(&category.ID)
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}
	return categories, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int) (*models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`

	var cat models.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&cat.ID, &cat.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}
