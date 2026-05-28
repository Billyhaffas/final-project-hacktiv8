package middleware

import (
	"api-gateway/helper"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, helper.Fail("UNAUTHORIZED", "missing or invalid authorization header"))
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		var claims Claims
		_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_EXPIRED", "token has expired"))
			}
			return c.JSON(http.StatusUnauthorized, helper.Fail("TOKEN_INVALID", "invalid token"))
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		return next(c)
	}
}
