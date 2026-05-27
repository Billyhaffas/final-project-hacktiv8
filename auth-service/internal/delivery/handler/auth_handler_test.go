package handler_test

import (
	"auth-service/internal/delivery/handler"
	"auth-service/internal/domain"
	"auth-service/internal/model/user"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
)

// --- mock usecase ---

type mockAuthUC struct {
	register       func(context.Context, user.RegisterRequest) (*user.RegisterResponse, error)
	login          func(context.Context, user.LoginRequest) (*user.LoginResponse, error)
	refresh        func(context.Context, user.RefreshRequest) (*user.RefreshResponse, error)
	logout         func(context.Context, user.LogoutRequest) error
	forgotPassword func(context.Context, user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error)
	resetPassword  func(context.Context, user.ResetPasswordRequest) error
}

func (m *mockAuthUC) Register(ctx context.Context, req user.RegisterRequest) (*user.RegisterResponse, error) {
	if m.register != nil {
		return m.register(ctx, req)
	}
	return nil, nil
}

func (m *mockAuthUC) Login(ctx context.Context, req user.LoginRequest) (*user.LoginResponse, error) {
	if m.login != nil {
		return m.login(ctx, req)
	}
	return nil, nil
}

func (m *mockAuthUC) Refresh(ctx context.Context, req user.RefreshRequest) (*user.RefreshResponse, error) {
	if m.refresh != nil {
		return m.refresh(ctx, req)
	}
	return nil, nil
}

func (m *mockAuthUC) Logout(ctx context.Context, req user.LogoutRequest) error {
	if m.logout != nil {
		return m.logout(ctx, req)
	}
	return nil
}

func (m *mockAuthUC) ForgotPassword(ctx context.Context, req user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error) {
	if m.forgotPassword != nil {
		return m.forgotPassword(ctx, req)
	}
	return &user.ForgotPasswordResponse{}, nil
}

func (m *mockAuthUC) ResetPassword(ctx context.Context, req user.ResetPasswordRequest) error {
	if m.resetPassword != nil {
		return m.resetPassword(ctx, req)
	}
	return nil
}

// newCtx builds an Echo context around a JSON body. e.NewContext already returns
// *echo.Context in v5.1.1, which matches the handler signature directly.
func newCtx(e *echo.Echo, body string) (*echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// --- Register handler ---

func TestRegisterHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		uc         *mockAuthUC
		wantStatus int
	}{
		{
			name: "201 — success",
			body: `{"name":"Alice","email":"alice@example.com","password":"secret123"}`,
			uc: &mockAuthUC{register: func(_ context.Context, _ user.RegisterRequest) (*user.RegisterResponse, error) {
				return &user.RegisterResponse{UserID: 1, Email: "alice@example.com", Name: "Alice"}, nil
			}},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "400 — malformed JSON",
			body:       `{invalid`,
			uc:         &mockAuthUC{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "400 — missing required fields",
			body: `{"email":"alice@example.com","password":"secret123"}`,
			uc: &mockAuthUC{register: func(_ context.Context, _ user.RegisterRequest) (*user.RegisterResponse, error) {
				return nil, domain.ErrRegisterFields
			}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "409 — email already registered",
			body: `{"name":"Alice","email":"alice@example.com","password":"secret123"}`,
			uc: &mockAuthUC{register: func(_ context.Context, _ user.RegisterRequest) (*user.RegisterResponse, error) {
				return nil, domain.ErrEmailTaken
			}},
			wantStatus: http.StatusConflict,
		},
		{
			name: "500 — unexpected repo error",
			body: `{"name":"Alice","email":"alice@example.com","password":"secret123"}`,
			uc: &mockAuthUC{register: func(_ context.Context, _ user.RegisterRequest) (*user.RegisterResponse, error) {
				return nil, errors.New("db failure")
			}},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAuthHandler(tt.uc)
			c, rec := newCtx(e, tt.body)
			if err := h.Register(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
			var body map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("could not parse response body: %v", err)
			}
			if rec.Code >= 400 {
				if body["success"] != false {
					t.Fatalf("want success=false on error, got %v", body["success"])
				}
				if body["error"] == nil {
					t.Fatal("expected error field on error response")
				}
			} else {
				if body["success"] != true {
					t.Fatalf("want success=true, got %v", body["success"])
				}
			}
		})
	}
}

// --- Login handler ---

func TestLoginHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		uc         *mockAuthUC
		wantStatus int
		checkToken bool
	}{
		{
			name: "200 — success with token in response",
			body: `{"email":"alice@example.com","password":"secret123"}`,
			uc: &mockAuthUC{login: func(_ context.Context, _ user.LoginRequest) (*user.LoginResponse, error) {
				return &user.LoginResponse{Token: "jwt.token.here", UserID: 1, Email: "alice@example.com"}, nil
			}},
			wantStatus: http.StatusOK,
			checkToken: true,
		},
		{
			name:       "400 — malformed JSON",
			body:       `{invalid`,
			uc:         &mockAuthUC{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "400 — missing fields",
			body: `{"email":"alice@example.com"}`,
			uc: &mockAuthUC{login: func(_ context.Context, _ user.LoginRequest) (*user.LoginResponse, error) {
				return nil, domain.ErrMissingFields
			}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "401 — invalid credentials",
			body: `{"email":"alice@example.com","password":"wrong"}`,
			uc: &mockAuthUC{login: func(_ context.Context, _ user.LoginRequest) (*user.LoginResponse, error) {
				return nil, domain.ErrInvalidCredentials
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "500 — unexpected error",
			body: `{"email":"alice@example.com","password":"secret123"}`,
			uc: &mockAuthUC{login: func(_ context.Context, _ user.LoginRequest) (*user.LoginResponse, error) {
				return nil, errors.New("db failure")
			}},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAuthHandler(tt.uc)
			c, rec := newCtx(e, tt.body)
			if err := h.Login(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
			var body map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("could not parse response body: %v", err)
			}
			if rec.Code >= 400 {
				if body["success"] != false {
					t.Fatalf("want success=false on error, got %v", body["success"])
				}
				if body["error"] == nil {
					t.Fatal("expected error field on error response")
				}
			}
			if tt.checkToken {
				if body["success"] != true {
					t.Fatalf("want success=true, got %v", body["success"])
				}
				data, ok := body["data"].(map[string]interface{})
				if !ok {
					t.Fatal("response data field is missing or not an object")
				}
				if data["token"] == "" || data["token"] == nil {
					t.Fatal("expected non-empty token in response data")
				}
			}
		})
	}
}

// --- ForgotPassword handler ---

func TestForgotPasswordHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		uc         *mockAuthUC
		wantStatus int
	}{
		{
			name: "200 — success",
			body: `{"email":"alice@example.com"}`,
			uc: &mockAuthUC{forgotPassword: func(_ context.Context, _ user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error) {
				return &user.ForgotPasswordResponse{ResetToken: "rawtoken"}, nil
			}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "400 — malformed JSON",
			body:       `{invalid`,
			uc:         &mockAuthUC{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "400 — missing email",
			body: `{}`,
			uc: &mockAuthUC{forgotPassword: func(_ context.Context, _ user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error) {
				return nil, domain.ErrMissingFields
			}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "500 — unexpected error",
			body: `{"email":"alice@example.com"}`,
			uc: &mockAuthUC{forgotPassword: func(_ context.Context, _ user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error) {
				return nil, errors.New("db failure")
			}},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAuthHandler(tt.uc)
			c, rec := newCtx(e, tt.body)
			if err := h.ForgotPassword(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
			var body map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("could not parse response body: %v", err)
			}
			if rec.Code >= 400 {
				if body["success"] != false {
					t.Fatalf("want success=false on error, got %v", body["success"])
				}
				if body["error"] == nil {
					t.Fatal("expected error field on error response")
				}
			} else {
				if body["success"] != true {
					t.Fatalf("want success=true, got %v", body["success"])
				}
			}
		})
	}
}

// --- ResetPassword handler ---

func TestResetPasswordHandler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		uc         *mockAuthUC
		wantStatus int
	}{
		{
			name: "200 — success",
			body: `{"token":"rawtoken","new_password":"newpass123"}`,
			uc:   &mockAuthUC{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "400 — malformed JSON",
			body:       `{invalid`,
			uc:         &mockAuthUC{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "400 — missing fields",
			body: `{"token":"rawtoken"}`,
			uc: &mockAuthUC{resetPassword: func(_ context.Context, _ user.ResetPasswordRequest) error {
				return domain.ErrMissingFields
			}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "401 — invalid token",
			body: `{"token":"bad","new_password":"newpass123"}`,
			uc: &mockAuthUC{resetPassword: func(_ context.Context, _ user.ResetPasswordRequest) error {
				return domain.ErrResetTokenInvalid
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "401 — token already used",
			body: `{"token":"used","new_password":"newpass123"}`,
			uc: &mockAuthUC{resetPassword: func(_ context.Context, _ user.ResetPasswordRequest) error {
				return domain.ErrResetTokenUsed
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "401 — token expired",
			body: `{"token":"old","new_password":"newpass123"}`,
			uc: &mockAuthUC{resetPassword: func(_ context.Context, _ user.ResetPasswordRequest) error {
				return domain.ErrResetTokenExpired
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "500 — unexpected error",
			body: `{"token":"rawtoken","new_password":"newpass123"}`,
			uc: &mockAuthUC{resetPassword: func(_ context.Context, _ user.ResetPasswordRequest) error {
				return errors.New("db failure")
			}},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAuthHandler(tt.uc)
			c, rec := newCtx(e, tt.body)
			if err := h.ResetPassword(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
			var body map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("could not parse response body: %v", err)
			}
			if rec.Code >= 400 {
				if body["success"] != false {
					t.Fatalf("want success=false on error, got %v", body["success"])
				}
				if body["error"] == nil {
					t.Fatal("expected error field on error response")
				}
			} else {
				if body["success"] != true {
					t.Fatalf("want success=true, got %v", body["success"])
				}
			}
		})
	}
}
