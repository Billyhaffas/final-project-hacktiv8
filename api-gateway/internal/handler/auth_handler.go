package handler

import (
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	proxy *httputil.ReverseProxy
}

func NewAuthHandler() *AuthHandler {
	target, _ := url.Parse(os.Getenv("AUTH_SERVICE_URL"))
	return &AuthHandler{proxy: httputil.NewSingleHostReverseProxy(target)}
}

// Register godoc
// @Summary      Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Registration data"
// @Success      201   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Login godoc
// @Summary      Login and get JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Login credentials"
// @Success      200   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshRequest  true  "Refresh token"
// @Success      200   {object}  SuccessResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Logout godoc
// @Summary      Revoke refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      LogoutRequest  true  "Refresh token to revoke"
// @Success      200   {object}  SuccessResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// ForgotPassword godoc
// @Summary      Request password reset email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      ForgotPasswordRequest  true  "Email address"
// @Success      200   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// ResetPassword godoc
// @Summary      Complete password reset
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      ResetPasswordRequest  true  "Reset token and new password"
// @Success      200   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *echo.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}
