package domain

import (
	"auth-service/internal/model/user"
	"context"
	"errors"
	"time"

	"github.com/labstack/echo/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrMissingFields      = errors.New("email and password are required")
	ErrEmailTaken         = errors.New("email already registered")
	ErrRegisterFields     = errors.New("name, email, and password are required")
	ErrTokenInvalid       = errors.New("refresh token is invalid")
	ErrTokenRevoked       = errors.New("refresh token has been revoked")
	ErrTokenExpired       = errors.New("refresh token has expired")
	ErrResetTokenInvalid  = errors.New("reset token is invalid")
	ErrResetTokenUsed     = errors.New("reset token has already been used")
	ErrResetTokenExpired  = errors.New("reset token has expired")
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUser(ctx context.Context, name, email, passwordHash string) (int, error)
	StoreRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*user.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	StorePasswordResetToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindPasswordResetToken(ctx context.Context, tokenHash string) (*user.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
}

type AuthUseCase interface {
	Register(ctx context.Context, req user.RegisterRequest) (*user.RegisterResponse, error)
	Login(ctx context.Context, req user.LoginRequest) (*user.LoginResponse, error)
	Refresh(ctx context.Context, req user.RefreshRequest) (*user.RefreshResponse, error)
	Logout(ctx context.Context, req user.LogoutRequest) error
	ForgotPassword(ctx context.Context, req user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req user.ResetPasswordRequest) error
}

type AuthHandler interface {
	Register(c *echo.Context) error
	Login(c *echo.Context) error
	Refresh(c *echo.Context) error
	Logout(c *echo.Context) error
	ForgotPassword(c *echo.Context) error
	ResetPassword(c *echo.Context) error
}
