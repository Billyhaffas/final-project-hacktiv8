package main

import (
	"auth-service/internal/delivery/handler"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	database.Connect()

	userRepo := repository.NewUserRepository(database.DB)
	authUC := usecase.NewAuthUseCase(userRepo)
	authHandler := handler.NewAuthHandler(authUC)

	e := echo.New()

	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/logout", authHandler.Logout)
	auth.POST("/forgot-password", authHandler.ForgotPassword)
	auth.POST("/reset-password", authHandler.ResetPassword)

	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}
	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("server stopped", "error", err)
	}
}
