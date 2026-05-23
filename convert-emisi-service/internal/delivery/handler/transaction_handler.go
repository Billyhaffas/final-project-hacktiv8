package handler

import (
	"net/http"
	"p3-lc01-billyhaffas/helper"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/transaction"

	"github.com/labstack/echo/v5"
)

type TransactionHandler struct {
	TransactionUseCase domain.TransactionUseCase
}

func NewTransactionHandler(TransactionUC domain.TransactionUseCase) domain.TransactionHandler {
	return &TransactionHandler{TransactionUseCase: TransactionUC}
}

func (uh *TransactionHandler) CreateTransaction(c *echo.Context) error {
	var requestTransaction *transaction.TransactionBodyPayload
	err := c.Bind(&requestTransaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.Respon{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
	}
	err = uh.TransactionUseCase.CreateTransaction(c.Request().Context(), requestTransaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, helper.Respon{
		Status:  http.StatusText(http.StatusCreated),
		Message: "Transaction has been created",
	})
}
