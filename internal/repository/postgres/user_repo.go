package postgres

import (
	"context"
	"errors"

	"shop/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepo implements repository.UserRepository for PostgreSQL.
type UserRepo struct {
	db *pgxpool.Pool
}

// NewUserRepository constructs a new UserRepo instance.
func NewUserRepository(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// Create inserts a new user record into database.
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByID fetches a user by primary key.
func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail fetches a user by unique email address.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update modifies email and password_hash of an existing user.
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, user.Email, user.PasswordHash, user.ID)
	return err
}

// UpdateRole modifies the role of a user.
func (r *UserRepo) UpdateRole(ctx context.Context, userID int64, role domain.Role) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, role, userID)
	return err
}
