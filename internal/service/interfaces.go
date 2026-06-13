package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByProvider(ctx context.Context, provider, providerID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetByToken(ctx context.Context, token string) (*uuid.UUID, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
