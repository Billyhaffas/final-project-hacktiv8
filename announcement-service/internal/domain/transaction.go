package domain

import (
	"context"
	"p3-lc01-billyhaffas/internal/model/transaction"

	"github.com/labstack/echo/v5"
	"github.com/robfig/cron/v3"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, room *transaction.TransactionInsertModel) error
	DeleteTransaction(ctx context.Context) error
}

type TransactionUseCase interface {
	DeleteTransaction(ctx context.Context) error
	CreateTransaction(ctx context.Context, request *transaction.TransactionBodyPayload) error
}

type TransactionHandler interface {
	CreateTransaction(c *echo.Context) error
}

type TransactionCron interface {
	DeleteCron(c *cron.Cron)
}
