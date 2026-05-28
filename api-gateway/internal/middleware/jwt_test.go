package middleware_test

import (
	"api-gateway/internal/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func makeToken(t *testing.T, secret string, userID int, email string, expiry time.Duration) string {
	t.Helper()
	claims := middleware.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("makeToken: %v", err)
	}
	return s
}

func TestJWTMiddleware(t *testing.T) {
	const secret = "test-secret"
	os.Setenv("JWT_SECRET", secret)
	t.Cleanup(func() { os.Unsetenv("JWT_SECRET") })

	e := echo.New()

	next := func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user_id": c.Get("user_id"),
			"email":   c.Get("email"),
		})
	}
	h := middleware.JWT(next)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "200 — valid token",
			authHeader: "Bearer " + makeToken(t, secret, 7, "alice@example.com", time.Hour),
			wantStatus: http.StatusOK,
		},
		{
			name:       "401 — missing header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "401 — no Bearer prefix",
			authHeader: makeToken(t, secret, 7, "alice@example.com", time.Hour),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "401 — expired token",
			authHeader: "Bearer " + makeToken(t, secret, 7, "alice@example.com", -time.Hour),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "401 — wrong secret",
			authHeader: "Bearer " + makeToken(t, "wrong-secret", 7, "alice@example.com", time.Hour),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "401 — malformed token",
			authHeader: "Bearer notavalidtoken",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = h(c)
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestJWTMiddleware_InjectsUserID(t *testing.T) {
	const secret = "test-secret"
	os.Setenv("JWT_SECRET", secret)
	t.Cleanup(func() { os.Unsetenv("JWT_SECRET") })

	e := echo.New()

	var gotUserID interface{}
	var gotEmail interface{}
	next := func(c *echo.Context) error {
		gotUserID = c.Get("user_id")
		gotEmail = c.Get("email")
		return c.JSON(http.StatusOK, nil)
	}
	h := middleware.JWT(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+makeToken(t, secret, 42, "bob@example.com", time.Hour))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = h(c)

	if gotUserID != 42 {
		t.Fatalf("want user_id=42, got %v", gotUserID)
	}
	if gotEmail != "bob@example.com" {
		t.Fatalf("want email=bob@example.com, got %v", gotEmail)
	}
}
