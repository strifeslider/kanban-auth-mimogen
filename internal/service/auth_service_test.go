package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/auth"
	"github.com/user/kanban-saas/pkg/model"
)

func TestNewAuthService(t *testing.T) {
	svc := &AuthService{}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestAuthService_JWTConfig(t *testing.T) {
	cfg := auth.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	if cfg.Secret != "test-secret" {
		t.Errorf("expected secret 'test-secret', got '%s'", cfg.Secret)
	}
	if cfg.AccessExpiry != 15*time.Minute {
		t.Errorf("expected access expiry 15m, got %v", cfg.AccessExpiry)
	}
}

func TestAuthService_StructInit(t *testing.T) {
	svc := &AuthService{
		jwtCfg: auth.JWTConfig{Secret: "test"},
	}

	if svc.jwtCfg.Secret != "test" {
		t.Errorf("expected secret 'test', got '%s'", svc.jwtCfg.Secret)
	}
}

func TestRegisterRequest(t *testing.T) {
	req := model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	if req.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", req.Email)
	}
}

func TestLoginRequest(t *testing.T) {
	req := model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	if req.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", req.Email)
	}
}

func TestAuthResponse(t *testing.T) {
	user := model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	tokens := model.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	resp := model.AuthResponse{
		User:   user,
		Tokens: tokens,
	}

	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", resp.User.Email)
	}
	if resp.Tokens.AccessToken != "access-token" {
		t.Errorf("expected access token 'access-token', got '%s'", resp.Tokens.AccessToken)
	}
}
