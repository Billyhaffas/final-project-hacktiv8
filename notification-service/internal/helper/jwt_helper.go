package helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// Helper method to parse out user context details safely from context values
func GetUserIDFromJWT(e *echo.Context) (int, error) {
	// Echo's JWT middleware binds the *jwt.Token object under the key "user" by default
	tokenVal := e.Get("user")
	if tokenVal == nil {
		return 0, errors.New("unauthorized: token data not found in context lifecycle")
	}

	token, ok := tokenVal.(*jwt.Token)
	if !ok {
		return 0, errors.New("unauthorized: invalid token instance layout")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("unauthorized: failed parsing token claims map")
	}

	// Extract the user_id property
	idVal, exists := claims["user_id"]
	if !exists {
		idVal, exists = claims["sub"] // Fallback checking standard 'sub' claim
		if !exists {
			return 0, errors.New("unauthorized: missing user identifier within token body")
		}
	}

	// Safely map json float64 numbers to integer types
	switch v := idVal.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, errors.New("unauthorized: corrupt user identity format")
	}
}
