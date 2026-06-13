package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, name, avatar_url, password_hash, provider, provider_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		user.ID, user.Email, user.Name, user.AvatarURL,
		user.PasswordHash, user.Provider, user.ProviderID,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, name, avatar_url, password_hash, provider, provider_id, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.AvatarURL,
		&user.PasswordHash, &user.Provider, &user.ProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, name, avatar_url, password_hash, provider, provider_id, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.AvatarURL,
		&user.PasswordHash, &user.Provider, &user.ProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByProvider(ctx context.Context, provider, providerID string) (*model.User, error) {
	query := `
		SELECT id, email, name, avatar_url, password_hash, provider, provider_id, created_at, updated_at
		FROM users WHERE provider = $1 AND provider_id = $2 AND deleted_at IS NULL`

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, provider, providerID).Scan(
		&user.ID, &user.Email, &user.Name, &user.AvatarURL,
		&user.PasswordHash, &user.Provider, &user.ProviderID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by provider: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET name = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, user.ID, user.Name, user.AvatarURL).Scan(&user.UpdatedAt)
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}
