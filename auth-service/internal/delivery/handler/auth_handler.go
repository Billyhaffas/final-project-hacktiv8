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
		return c.JSON(http.StatusBadRequest, helper.Response{
			Status:  "Bad Request",
			Message: "invalid request body",
		})
	}

	resp, err := h.authUC.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrRegisterFields) {
			return c.JSON(http.StatusBadRequest, helper.Response{
				Status:  "Bad Request",
				Message: err.Error(),
			})
		}
		if errors.Is(err, domain.ErrEmailTaken) {
			return c.JSON(http.StatusConflict, helper.Response{
				Status:  "Conflict",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, helper.Response{
			Status:  "Internal Server Error",
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, helper.Response{
		Status:  "Created",
		Message: "user registered successfully",
		Data:    resp,
	})
}

func (h *authHandler) Login(c *echo.Context) error {
	var req user.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Response{
			Status:  "Bad Request",
			Message: "invalid request body",
		})
	}

	resp, err := h.authUC.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrMissingFields) {
			return c.JSON(http.StatusBadRequest, helper.Response{
				Status:  "Bad Request",
				Message: err.Error(),
			})
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, helper.Response{
				Status:  "Unauthorized",
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, helper.Response{
			Status:  "Internal Server Error",
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, helper.Response{
		Status:  "OK",
		Message: "login successful",
		Data:    resp,
	})
}
