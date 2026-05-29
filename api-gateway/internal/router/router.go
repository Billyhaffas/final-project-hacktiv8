package router

import (
	"api-gateway/internal/handler"
	mw "api-gateway/internal/middleware"

	_ "api-gateway/docs"

	"github.com/labstack/echo/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/labstack/echo/v5"
)

func Setup(e *echo.Echo, auth *handler.AuthHandler, emission *handler.EmissionHandler,
	pref *handler.PreferenceHandler, alert *handler.AlertHandler, cvt *handler.ConvertHandler) {

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Swagger UI — http://host:8080/swagger/index.html
	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.WrapHandler))

	api := e.Group("/api/v1")

	// Auth — proxied to auth-service
	authGroup := api.Group("/auth")
	authGroup.POST("/register", auth.Register)
	authGroup.POST("/login", auth.Login)
	authGroup.POST("/refresh", auth.Refresh)
	authGroup.POST("/forgot-password", auth.ForgotPassword)
	authGroup.POST("/reset-password", auth.ResetPassword)
	authGroup.POST("/logout", auth.Logout, mw.JWT)

	// Emissions — all require JWT
	emissionGroup := api.Group("/emissions", mw.JWT)
	emissionGroup.POST("", emission.LogEmission)
	emissionGroup.GET("/today", emission.GetDailyTotal)
	emissionGroup.GET("/alert", alert.CheckAlert)
	emissionGroup.GET("/report", emission.GetMonthlyReport)
	emissionGroup.GET("/convert", cvt.ConvertToIDR)

	// Preferences — require JWT
	prefGroup := api.Group("/preferences", mw.JWT)
	prefGroup.GET("", pref.GetPreferences)
	prefGroup.PUT("", pref.UpdatePreferences)
}
