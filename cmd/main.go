package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/auth"
	"github.com/user/kanban-saas/pkg/database"
	appmiddleware "github.com/user/kanban-saas/pkg/middleware"
	"github.com/user/kanban-saas/services/auth/internal/handler"
	"github.com/user/kanban-saas/services/auth/internal/repository"
	"github.com/user/kanban-saas/services/auth/internal/service"
)

func main() {
	env := getEnv("ENV", "local")
	port := getEnv("PORT", "8081")

	logger := setupLogger(env)
	logger.Info("starting auth service", "env", env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewPostgresPool(ctx, database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "kanban"),
		Password: getEnv("DB_PASSWORD", "kanban_dev_password"),
		Database: getEnv("DB_NAME", "kanban_auth"),
		MaxConns: 10,
		MinConns: 2,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	runMigrations(ctx, db, logger)

	jwtCfg := auth.JWTConfig{
		Secret:        getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 168 * time.Hour,
	}

	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)
	authService := service.NewAuthService(userRepo, refreshRepo, jwtCfg)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()

	allowedOrigins := appmiddleware.ParseOrigins(getEnv("CORS_ORIGINS", "http://localhost:3000"))
	r.Use(appmiddleware.CORS(allowedOrigins))
	r.Use(appmiddleware.Logging(logger))
	r.Use(appmiddleware.Recovery(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler.SetupRoutes(r, authHandler, jwtCfg)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("auth service listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down auth service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
	logger.Info("auth service stopped")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {
	case "local":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(handler)
}

func runMigrations(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			avatar_url TEXT,
			password_hash VARCHAR(255),
			provider VARCHAR(50) NOT NULL DEFAULT 'local',
			provider_id VARCHAR(255),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(500) NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			revoked_at TIMESTAMPTZ
		);`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(ctx, m); err != nil {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}
	}
	logger.Info("migrations completed")
}
