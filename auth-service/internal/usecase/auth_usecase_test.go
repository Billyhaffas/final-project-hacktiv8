package usecase_test

import (
	"auth-service/internal/domain"
	"auth-service/internal/model/user"
	"auth-service/internal/usecase"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- mock repository ---

type mockRepo struct {
	findByEmail func(context.Context, string) (*user.User, error)
	createUser  func(context.Context, string, string, string) (int, error)
}

func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.findByEmail != nil {
		return m.findByEmail(ctx, email)
	}
	return nil, nil
}

func (m *mockRepo) CreateUser(ctx context.Context, name, email, hash string) (int, error) {
	if m.createUser != nil {
		return m.createUser(ctx, name, email, hash)
	}
	return 0, nil
}

// bcryptHash produces a bcrypt hash using MinCost so tests run fast.
func bcryptHash(t *testing.T, pw string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcryptHash: %v", err)
	}
	return string(h)
}

// --- Register ---

func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		req     user.RegisterRequest
		repo    *mockRepo
		wantErr error
		wantID  int
	}{
		{
			name: "success",
			req:  user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret123"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) { return nil, nil },
				createUser:  func(_ context.Context, _, _, _ string) (int, error) { return 42, nil },
			},
			wantID: 42,
		},
		{
			name:    "missing name",
			req:     user.RegisterRequest{Email: "alice@example.com", Password: "secret123"},
			repo:    &mockRepo{},
			wantErr: domain.ErrRegisterFields,
		},
		{
			name:    "missing email",
			req:     user.RegisterRequest{Name: "Alice", Password: "secret123"},
			repo:    &mockRepo{},
			wantErr: domain.ErrRegisterFields,
		},
		{
			name:    "missing password",
			req:     user.RegisterRequest{Name: "Alice", Email: "alice@example.com"},
			repo:    &mockRepo{},
			wantErr: domain.ErrRegisterFields,
		},
		{
			name: "email already taken",
			req:  user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret123"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) {
					return &user.User{ID: 1, Email: "alice@example.com"}, nil
				},
			},
			wantErr: domain.ErrEmailTaken,
		},
		{
			name: "repo error on FindByEmail",
			req:  user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret123"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) {
					return nil, errors.New("db unavailable")
				},
			},
			wantErr: errors.New("any wrapped error"),
		},
		{
			name: "repo error on CreateUser",
			req:  user.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "secret123"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) { return nil, nil },
				createUser: func(_ context.Context, _, _, _ string) (int, error) {
					return 0, errors.New("insert failed")
				},
			},
			wantErr: errors.New("any wrapped error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.NewAuthUseCase(tt.repo)
			resp, err := uc.Register(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				// sentinel errors must match exactly via errors.Is
				if errors.Is(tt.wantErr, domain.ErrRegisterFields) || errors.Is(tt.wantErr, domain.ErrEmailTaken) {
					if !errors.Is(err, tt.wantErr) {
						t.Fatalf("want %v, got %v", tt.wantErr, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.UserID != tt.wantID {
				t.Fatalf("want UserID=%d, got %d", tt.wantID, resp.UserID)
			}
			if resp.Email != tt.req.Email {
				t.Fatalf("want Email=%s, got %s", tt.req.Email, resp.Email)
			}
			if resp.Name != tt.req.Name {
				t.Fatalf("want Name=%s, got %s", tt.req.Name, resp.Name)
			}
		})
	}
}

// --- Login ---

func TestLogin(t *testing.T) {
	const secret = "test-jwt-secret"
	os.Setenv("JWT_SECRET", secret)
	t.Cleanup(func() { os.Unsetenv("JWT_SECRET") })

	validHash := bcryptHash(t, "correct-pass")

	tests := []struct {
		name       string
		req        user.LoginRequest
		repo       *mockRepo
		wantErr    error
		checkToken bool
	}{
		{
			name: "success — valid token with correct claims",
			req:  user.LoginRequest{Email: "alice@example.com", Password: "correct-pass"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) {
					return &user.User{ID: 7, Email: "alice@example.com", PasswordHash: validHash}, nil
				},
			},
			checkToken: true,
		},
		{
			name:    "missing email",
			req:     user.LoginRequest{Password: "correct-pass"},
			repo:    &mockRepo{},
			wantErr: domain.ErrMissingFields,
		},
		{
			name:    "missing password",
			req:     user.LoginRequest{Email: "alice@example.com"},
			repo:    &mockRepo{},
			wantErr: domain.ErrMissingFields,
		},
		{
			name: "user not found",
			req:  user.LoginRequest{Email: "ghost@example.com", Password: "correct-pass"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) { return nil, nil },
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			req:  user.LoginRequest{Email: "alice@example.com", Password: "wrong"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) {
					return &user.User{ID: 7, Email: "alice@example.com", PasswordHash: validHash}, nil
				},
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "repo error",
			req:  user.LoginRequest{Email: "alice@example.com", Password: "correct-pass"},
			repo: &mockRepo{
				findByEmail: func(_ context.Context, _ string) (*user.User, error) {
					return nil, errors.New("connection timeout")
				},
			},
			wantErr: errors.New("any wrapped error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.NewAuthUseCase(tt.repo)
			resp, err := uc.Login(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tt.wantErr, domain.ErrMissingFields) || errors.Is(tt.wantErr, domain.ErrInvalidCredentials) {
					if !errors.Is(err, tt.wantErr) {
						t.Fatalf("want %v, got %v", tt.wantErr, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected response, got nil")
			}

			if tt.checkToken {
				var claims struct {
					UserID int    `json:"user_id"`
					Email  string `json:"email"`
					jwt.RegisteredClaims
				}
				parsed, parseErr := jwt.ParseWithClaims(resp.Token, &claims, func(tok *jwt.Token) (interface{}, error) {
					if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
						t.Fatalf("unexpected signing method: %v", tok.Header["alg"])
					}
					return []byte(secret), nil
				})
				if parseErr != nil || !parsed.Valid {
					t.Fatalf("invalid JWT: %v", parseErr)
				}
				if claims.UserID != 7 {
					t.Fatalf("want user_id=7, got %d", claims.UserID)
				}
				if claims.Email != "alice@example.com" {
					t.Fatalf("want email=alice@example.com, got %s", claims.Email)
				}
				exp := claims.ExpiresAt.Time
				if exp.Before(time.Now().Add(23*time.Hour)) || exp.After(time.Now().Add(25*time.Hour)) {
					t.Fatalf("expiry out of expected 24h window: %v", exp)
				}
				if resp.UserID != 7 {
					t.Fatalf("want resp.UserID=7, got %d", resp.UserID)
				}
				if resp.Email != "alice@example.com" {
					t.Fatalf("want resp.Email=alice@example.com, got %s", resp.Email)
				}
			}
		})
	}
}
