package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/auth"
	"github.com/user/kanban-saas/pkg/mock"
	"github.com/user/kanban-saas/pkg/model"
)

func newTestAuthService() *AuthService {
	userRepo := mock.NewMockUserRepo()
	refreshRepo := mock.NewMockRefreshTokenRepo()
	jwtCfg := auth.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}
	return &AuthService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtCfg:      jwtCfg,
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, err := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.User.Email)
	}
	if resp.Tokens.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if resp.Tokens.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "User 1",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "User 2",
		Password: "password456",
	})
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	resp, err := svc.Login(ctx, model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", resp.User.Email)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	_, err := svc.Login(ctx, model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestAuthService_Login_NonExistentEmail(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.Login(ctx, model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Error("expected error for nonexistent email")
	}
}

func TestAuthService_Login_SocialAccount(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	svc.userRepo.Create(ctx, &model.User{
		ID:       uuid.New(),
		Email:    "social@example.com",
		Name:     "Social",
		Provider: "google",
	})

	_, err := svc.Login(ctx, model.LoginRequest{
		Email:    "social@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Error("expected error for social account")
	}
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, _ := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	newResp, err := svc.RefreshToken(ctx, resp.Tokens.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newResp.Tokens.AccessToken == "" {
		t.Error("expected new access token")
	}
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.RefreshToken(ctx, "invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestAuthService_RefreshToken_ReuseToken(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, _ := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	svc.RefreshToken(ctx, resp.Tokens.RefreshToken)

	_, err := svc.RefreshToken(ctx, resp.Tokens.RefreshToken)
	if err == nil {
		t.Error("expected error for reused token")
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, _ := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	err := svc.Logout(ctx, resp.User.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	token := resp.Tokens.RefreshToken
	userID, _ := svc.refreshRepo.GetByToken(ctx, token)
	if userID != nil {
		t.Error("token should be revoked")
	}
}

func TestAuthService_GetProfile_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, _ := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	user, err := svc.GetProfile(ctx, resp.User.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}
}

func TestAuthService_GetProfile_NotFound(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.GetProfile(ctx, uuid.New())
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestAuthService_UpdateProfile_Success(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, _ := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	user, err := svc.UpdateProfile(ctx, resp.User.ID, "New Name", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "New Name" {
		t.Errorf("expected name New Name, got %s", user.Name)
	}
}

func TestAuthService_UpdateProfile_WithAvatar(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	resp, _ := svc.Register(ctx, model.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test",
		Password: "password123",
	})

	avatar := "https://example.com/avatar.png"
	user, err := svc.UpdateProfile(ctx, resp.User.ID, "", &avatar)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.AvatarURL == nil || *user.AvatarURL != avatar {
		t.Error("avatar not updated")
	}
}

func TestAuthService_UpdateProfile_NotFound(t *testing.T) {
	svc := newTestAuthService()
	ctx := context.Background()

	_, err := svc.UpdateProfile(ctx, uuid.New(), "Name", nil)
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}
