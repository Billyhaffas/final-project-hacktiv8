package usecase

import (
	"auth-service/internal/domain"
	"auth-service/internal/model/user"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const refreshTokenTTL = 7 * 24 * time.Hour

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

	return &user.RegisterResponse{UserID: id, Email: req.Email, Name: req.Name}, nil
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

	rawRefresh, hashRefresh, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("Login: generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(refreshTokenTTL)
	if err := uc.userRepo.StoreRefreshToken(ctx, u.ID, hashRefresh, expiresAt); err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	return &user.LoginResponse{
		Token:        token,
		RefreshToken: rawRefresh,
		UserID:       u.ID,
		Email:        u.Email,
	}, nil
}

func (uc *authUseCase) Refresh(ctx context.Context, req user.RefreshRequest) (*user.RefreshResponse, error) {
	if req.RefreshToken == "" {
		return nil, domain.ErrTokenInvalid
	}

	tokenHash := hashToken(req.RefreshToken)

	rt, err := uc.userRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}
	if rt == nil {
		return nil, domain.ErrTokenInvalid
	}
	if rt.IsRevoked {
		return nil, domain.ErrTokenRevoked
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	// rotate: revoke old token
	if err := uc.userRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("Refresh: revoke old token: %w", err)
	}

	newJWT, err := issueJWT(rt.UserID, rt.Email)
	if err != nil {
		return nil, fmt.Errorf("Refresh: issue JWT: %w", err)
	}

	rawRefresh, hashRefresh, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("Refresh: generate token: %w", err)
	}

	if err := uc.userRepo.StoreRefreshToken(ctx, rt.UserID, hashRefresh, time.Now().Add(refreshTokenTTL)); err != nil {
		return nil, fmt.Errorf("Refresh: store token: %w", err)
	}

	return &user.RefreshResponse{Token: newJWT, RefreshToken: rawRefresh}, nil
}

func (uc *authUseCase) Logout(ctx context.Context, req user.LogoutRequest) error {
	if req.RefreshToken == "" {
		return domain.ErrTokenInvalid
	}

	tokenHash := hashToken(req.RefreshToken)

	rt, err := uc.userRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("Logout: %w", err)
	}
	if rt == nil {
		return domain.ErrTokenInvalid
	}

	if err := uc.userRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("Logout: %w", err)
	}
	return nil
}

const resetTokenTTL = 1 * time.Hour

func (uc *authUseCase) ForgotPassword(ctx context.Context, req user.ForgotPasswordRequest) (*user.ForgotPasswordResponse, error) {
	if req.Email == "" {
		return nil, domain.ErrMissingFields
	}

	u, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("ForgotPassword: %w", err)
	}
	if u == nil {
		// Don't reveal whether the email exists — return success with empty token
		return &user.ForgotPasswordResponse{ResetToken: ""}, nil
	}

	rawToken, tokenHash, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("ForgotPassword: generate token: %w", err)
	}

	expiresAt := time.Now().Add(resetTokenTTL)
	if err := uc.userRepo.StorePasswordResetToken(ctx, u.ID, tokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("ForgotPassword: %w", err)
	}

	return &user.ForgotPasswordResponse{ResetToken: rawToken}, nil
}

func (uc *authUseCase) ResetPassword(ctx context.Context, req user.ResetPasswordRequest) error {
	if req.Token == "" || req.NewPassword == "" {
		return domain.ErrMissingFields
	}

	th := hashToken(req.Token)

	rt, err := uc.userRepo.FindPasswordResetToken(ctx, th)
	if err != nil {
		return fmt.Errorf("ResetPassword: %w", err)
	}
	if rt == nil {
		return domain.ErrResetTokenInvalid
	}
	if rt.IsUsed {
		return domain.ErrResetTokenUsed
	}
	if time.Now().After(rt.ExpiresAt) {
		return domain.ErrResetTokenExpired
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ResetPassword: hash password: %w", err)
	}

	if err := uc.userRepo.UpdatePassword(ctx, rt.UserID, string(hash)); err != nil {
		return fmt.Errorf("ResetPassword: %w", err)
	}

	if err := uc.userRepo.MarkPasswordResetTokenUsed(ctx, th); err != nil {
		return fmt.Errorf("ResetPassword: mark token: %w", err)
	}

	return nil
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

func generateRefreshToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = hex.EncodeToString(b)
	hash = hashToken(raw)
	return
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
