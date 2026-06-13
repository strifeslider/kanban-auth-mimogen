package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/user/kanban-saas/pkg/auth"
	apperr "github.com/user/kanban-saas/pkg/errors"
	"github.com/user/kanban-saas/pkg/model"
	"github.com/user/kanban-saas/services/auth/internal/repository"
)

type AuthService struct {
	userRepo       *repository.UserRepository
	refreshRepo    *repository.RefreshTokenRepository
	jwtCfg         auth.JWTConfig
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshRepo *repository.RefreshTokenRepository,
	jwtCfg auth.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		jwtCfg:      jwtCfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, apperr.Internal("failed to check email")
	}
	if exists {
		return nil, apperr.Conflict("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperr.Internal("failed to hash password")
	}

	hashStr := string(hash)
	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: &hashStr,
		Provider:     "local",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperr.Internal("failed to create user")
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperr.Unauthorized("invalid email or password")
	}

	if user.PasswordHash == nil {
		return nil, apperr.Unauthorized("account uses social login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperr.Unauthorized("invalid email or password")
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	userID, err := s.refreshRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, apperr.Unauthorized("invalid refresh token")
	}

	if err := s.refreshRepo.Revoke(ctx, refreshToken); err != nil {
		return nil, apperr.Internal("failed to revoke token")
	}

	user, err := s.userRepo.GetByID(ctx, *userID)
	if err != nil {
		return nil, apperr.Unauthorized("user not found")
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.refreshRepo.RevokeAllForUser(ctx, userID)
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperr.NotFound("user not found")
	}
	return user, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, name string, avatarURL *string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperr.NotFound("user not found")
	}

	if name != "" {
		user.Name = name
	}
	user.AvatarURL = avatarURL

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperr.Internal("failed to update user")
	}

	return user, nil
}

func (s *AuthService) generateTokens(ctx context.Context, user *model.User) (*model.AuthResponse, error) {
	accessToken, err := auth.GenerateAccessToken(s.jwtCfg, user.ID, user.Email)
	if err != nil {
		return nil, apperr.Internal("failed to generate access token")
	}

	refreshToken, err := auth.GenerateRefreshToken(s.jwtCfg)
	if err != nil {
		return nil, apperr.Internal("failed to generate refresh token")
	}

	expiresAt := time.Now().Add(s.jwtCfg.RefreshExpiry)
	if err := s.refreshRepo.Create(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, apperr.Internal("failed to store refresh token")
	}

	return &model.AuthResponse{
		User: *user,
		Tokens: model.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
