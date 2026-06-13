package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, userID, token, expiresAt)
	return err
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*uuid.UUID, error) {
	query := `
		SELECT user_id FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL`

	var userID uuid.UUID
	err := r.db.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &userID, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = $1`
	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
