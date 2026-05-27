package repository

import (
	"auth-service/internal/domain"
	"auth-service/internal/model/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, name FROM auth_db.users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindByEmail: %w", err)
	}
	return &u, nil
}

func (r *userRepository) CreateUser(ctx context.Context, name, email, passwordHash string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO auth_db.users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		name, email, passwordHash,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateUser: %w", err)
	}
	return id, nil
}

func (r *userRepository) StoreRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO auth_db.refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("StoreRefreshToken: %w", err)
	}
	return nil
}

func (r *userRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*user.RefreshToken, error) {
	var rt user.RefreshToken
	err := r.db.QueryRowContext(ctx,
		`SELECT rt.user_id, rt.token_hash, rt.expires_at, rt.is_revoked, u.email
		 FROM auth_db.refresh_tokens rt
		 JOIN auth_db.users u ON u.id = rt.user_id
		 WHERE rt.token_hash = $1`,
		tokenHash,
	).Scan(&rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.IsRevoked, &rt.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindRefreshToken: %w", err)
	}
	return &rt, nil
}

func (r *userRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE auth_db.refresh_tokens SET is_revoked = TRUE WHERE token_hash = $1`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("RevokeRefreshToken: %w", err)
	}
	return nil
}

func (r *userRepository) StorePasswordResetToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO auth_db.password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("StorePasswordResetToken: %w", err)
	}
	return nil
}

func (r *userRepository) FindPasswordResetToken(ctx context.Context, tokenHash string) (*user.PasswordResetToken, error) {
	var rt user.PasswordResetToken
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, token_hash, expires_at, is_used FROM auth_db.password_reset_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.IsUsed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindPasswordResetToken: %w", err)
	}
	return &rt, nil
}

func (r *userRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE auth_db.password_reset_tokens SET is_used = TRUE WHERE token_hash = $1`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("MarkPasswordResetTokenUsed: %w", err)
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE auth_db.users SET password_hash = $1 WHERE id = $2`,
		passwordHash, userID,
	)
	if err != nil {
		return fmt.Errorf("UpdatePassword: %w", err)
	}
	return nil
}
