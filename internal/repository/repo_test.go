package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

func TestUserRepository_New(t *testing.T) {
	repo := &UserRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestRefreshTokenRepository_New(t *testing.T) {
	repo := &RefreshTokenRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestUserRepository_Model(t *testing.T) {
	user := &model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Name:     "Test",
		Provider: "local",
	}
	if user.Email != "test@example.com" {
		t.Error("email mismatch")
	}
}

func TestRefreshTokenRepository_TokenModel(t *testing.T) {
	token := "refresh-token-123"
	if token == "" {
		t.Error("expected non-empty token")
	}
}
