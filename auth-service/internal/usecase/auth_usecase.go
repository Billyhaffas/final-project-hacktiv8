package usecase

import (
	"auth-service/internal/domain"
	"auth-service/internal/model/user"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo domain.UserRepository
}

func NewAuthUseCase(userRepo domain.UserRepository) domain.AuthUseCase {
	return &authUseCase{userRepo: userRepo}
}

type jwtClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (uc *authUseCase) Login(ctx context.Context, req user.LoginRequest) (*user.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrMissingFields
	}

	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}
	if u == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := issueJWT(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("Login: issue token: %w", err)
	}

	return &user.LoginResponse{
		Token:  token,
		UserID: u.ID,
		Email:  u.Email,
	}, nil
}

func (uc *authUseCase) Register(ctx context.Context, req user.RegisterRequest) (*user.RegisterResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, domain.ErrRegisterFields
	}

	existing, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Register: hash password: %w", err)
	}

	id, err := uc.userRepo.CreateUser(ctx, req.Name, req.Email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	return &user.RegisterResponse{
		UserID: id,
		Email:  req.Email,
		Name:   req.Name,
	}, nil
}

func issueJWT(userID int, email string) (string, error) {
	claims := jwtClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
