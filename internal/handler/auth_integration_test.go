package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func TestAuthHandler_Register_EmptyBody(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("POST", "/api/v1/auth/register", nil)
	w := httptest.NewRecorder()
	h.Register(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString("{bad"))
	w := httptest.NewRecorder()
	h.Register(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	h := &AuthHandler{}
	body := map[string]string{"email": "", "name": "", "password": ""}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Register(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_EmptyBody(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	h := &AuthHandler{}
	body := map[string]string{"email": "", "password": ""}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Login(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_EmptyBody(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
	w := httptest.NewRecorder()
	h.RefreshToken(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_RefreshToken_EmptyToken(t *testing.T) {
	h := &AuthHandler{}
	body := map[string]string{"refresh_token": ""}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RefreshToken(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_GetProfile_NoAuth(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	h.GetProfile(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Logout_NoAuth(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()
	h.Logout(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateProfile_NoAuth(t *testing.T) {
	h := &AuthHandler{}
	body := map[string]string{"name": "New Name"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/auth/me", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateProfile(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_UpdateProfile_EmptyBody(t *testing.T) {
	h := &AuthHandler{}
	req := httptest.NewRequest("PUT", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	h.UpdateProfile(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSetupRoutes_Auth(t *testing.T) {
	r := chi.NewRouter()
	h := &AuthHandler{}
	jwtCfg := auth.JWTConfig{Secret: "test"}
	SetupRoutes(r, h, jwtCfg)

	routes := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
	}
	for _, route := range routes {
		req := httptest.NewRequest("POST", route, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusNotFound {
			t.Errorf("route %s not found", route)
		}
	}
}
