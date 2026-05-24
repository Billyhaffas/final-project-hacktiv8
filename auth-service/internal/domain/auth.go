package domain

import (
	"auth-service/internal/model/user"
	"context"
	"errors"

	"github.com/labstack/echo/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrMissingFields      = errors.New("email and password are required")
	ErrEmailTaken         = errors.New("email already registered")
	ErrRegisterFields     = errors.New("name, email, and password are required")
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUser(ctx context.Context, name, email, passwordHash string) (int, error)
}

type AuthUseCase interface {
	Login(ctx context.Context, req user.LoginRequest) (*user.LoginResponse, error)
	Register(ctx context.Context, req user.RegisterRequest) (*user.RegisterResponse, error)
}

type AuthHandler interface {
	Login(c *echo.Context) error
	Register(c *echo.Context) error
}
