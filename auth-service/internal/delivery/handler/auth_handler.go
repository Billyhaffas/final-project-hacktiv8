package handler

import (
	"auth-service/helper"
	"auth-service/internal/domain"
	"auth-service/internal/model/user"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type authHandler struct {
	authUC domain.AuthUseCase
}

func NewAuthHandler(authUC domain.AuthUseCase) domain.AuthHandler {
	return &authHandler{authUC: authUC}
}

func (h *authHandler) Register(c *echo.Context) error {
	var req user.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	resp, err := h.authUC.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrRegisterFields) {
			return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", err.Error()))
		}
		if errors.Is(err, domain.ErrEmailTaken) {
			return c.JSON(http.StatusConflict, helper.Fail("EMAIL_TAKEN", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "internal server error"))
	}

	return c.JSON(http.StatusCreated, helper.Success(resp))
}

func (h *authHandler) Login(c *echo.Context) error {
	var req user.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	resp, err := h.authUC.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrMissingFields) {
			return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", err.Error()))
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("INVALID_CREDENTIALS", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "internal server error"))
	}

	return c.JSON(http.StatusOK, helper.Success(resp))
}

func (h *authHandler) Refresh(c *echo.Context) error {
	var req user.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	resp, err := h.authUC.Refresh(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_INVALID", err.Error()))
		}
		if errors.Is(err, domain.ErrTokenRevoked) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_REVOKED", err.Error()))
		}
		if errors.Is(err, domain.ErrTokenExpired) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_EXPIRED", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "internal server error"))
	}

	return c.JSON(http.StatusOK, helper.Success(resp))
}

func (h *authHandler) ForgotPassword(c *echo.Context) error {
	var req user.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	resp, err := h.authUC.ForgotPassword(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrMissingFields) {
			return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "internal server error"))
	}

	return c.JSON(http.StatusOK, helper.Success(resp))
}

func (h *authHandler) ResetPassword(c *echo.Context) error {
	var req user.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	if err := h.authUC.ResetPassword(c.Request().Context(), req); err != nil {
		if errors.Is(err, domain.ErrMissingFields) {
			return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", err.Error()))
		}
		if errors.Is(err, domain.ErrResetTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_INVALID", err.Error()))
		}
		if errors.Is(err, domain.ErrResetTokenUsed) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_USED", err.Error()))
		}
		if errors.Is(err, domain.ErrResetTokenExpired) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_EXPIRED", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "internal server error"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]string{"message": "password reset successfully"}))
}

func (h *authHandler) Logout(c *echo.Context) error {
	var req user.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	if err := h.authUC.Logout(c.Request().Context(), req); err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_INVALID", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "internal server error"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]string{"message": "logged out successfully"}))
}
