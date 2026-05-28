package router

import (
	"api-gateway/internal/handler"
	mw "api-gateway/internal/middleware"

	"github.com/labstack/echo/v5"
)

func Setup(e *echo.Echo, auth *handler.AuthHandler, emission *handler.EmissionHandler, pref *handler.PreferenceHandler) {
	api := e.Group("/api/v1")

	// Auth — proxied to auth-service
	authGroup := api.Group("/auth")
	authGroup.POST("/register", auth.Proxy)
	authGroup.POST("/login", auth.Proxy)
	authGroup.POST("/refresh", auth.Proxy)
	authGroup.POST("/forgot-password", auth.Proxy)
	authGroup.POST("/reset-password", auth.Proxy)
	authGroup.POST("/logout", auth.Proxy, mw.JWT)

	// Emissions — all require JWT
	emissionGroup := api.Group("/emissions", mw.JWT)
	emissionGroup.POST("", emission.LogEmission)
	emissionGroup.GET("/today", emission.GetDailyTotal)
	emissionGroup.GET("/alert", emission.GetAlert)
	emissionGroup.GET("/report", emission.GetMonthlyReport)
	emissionGroup.GET("/convert", emission.ConvertToIDR)

	// Preferences — require JWT
	prefGroup := api.Group("/preferences", mw.JWT)
	prefGroup.GET("", pref.GetPreferences)
	prefGroup.PUT("", pref.UpdatePreferences)
}
