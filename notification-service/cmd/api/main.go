package main

import (
	"p3-lc01-billyhaffas/internal/delivery/cronjob"
	"p3-lc01-billyhaffas/internal/delivery/handler"
	"p3-lc01-billyhaffas/internal/infrastructure/database"
	"p3-lc01-billyhaffas/internal/repository"
	"p3-lc01-billyhaffas/internal/usecase"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/labstack/echo/v5"
)

func main() {
	database.ConnectDB()
	usersCollection := database.DB.Collection("Users")
	roomCollection := database.DB.Collection("Room")
	transactionCollection := database.DB.Collection("Transactions")

	//init repository
	userRepository := repository.NewUserCollection(usersCollection)
	roomRepository := repository.NewRoomCollection(roomCollection)
	transactionRepository := repository.NewTransactionCollection(transactionCollection)

	//init usecase
	userUseCase := usecase.NewUserUseCase(userRepository)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepository, userRepository, roomRepository)

	//init handler
	userHandler := handler.NewUserHandler(userUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)

	// init cron
	transactionCron := cronjob.NewTransactionCronjob(transactionUseCase)

	echo := echo.New()
	echo.POST("/users", userHandler.PostUser)
	echo.GET("/users/:id", userHandler.GetUserById)
	echo.GET("/users", userHandler.GetAllUser)
	echo.DELETE("/users/:id", userHandler.DeleteUserById)
	echo.PUT("/users/:id", userHandler.UpdateUserById)

	loc, _ := time.LoadLocation("Asia/Jakarta")

	cronJob := cron.New(cron.WithLocation(loc))
	transactionCron.DeleteCron(cronJob)
	cronJob.Start()

	echo.POST("/transaction", transactionHandler.CreateTransaction)

	if err := echo.Start(":" + "8085"); err != nil {
		echo.Logger.Error("failed to start server", "error", err)
	}
}
