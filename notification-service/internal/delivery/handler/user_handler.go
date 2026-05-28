package handler

import (
	"net/http"
	"p3-lc01-billyhaffas/helper"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/user"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	UserUseCase domain.UserUseCase
}

func NewUserHandler(userUC domain.UserUseCase) domain.UserHandler {
	return &UserHandler{UserUseCase: userUC}
}

func (uh *UserHandler) PostUser(c *echo.Context) error {
	var requestUser *user.UserInsertModel
	err := c.Bind(&requestUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.Respon{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
	}
	err = uh.UserUseCase.PostUser(c.Request().Context(), requestUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, helper.Respon{
		Status:  http.StatusText(http.StatusCreated),
		Message: "User has been created",
	})
}
func (uh *UserHandler) GetUserById(c *echo.Context) error {
	id := c.Param("id")
	respon, err := uh.UserUseCase.GetUserById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, helper.GetUserTypeRespon{
		Status:  http.StatusText(http.StatusOK),
		Message: "User successfully take it",
		Data:    respon,
	})
}

func (uh *UserHandler) GetAllUser(c *echo.Context) error {
	respon, err := uh.UserUseCase.GetAllUser(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, helper.GetAllUserTypeRespon{
		Status:  http.StatusText(http.StatusOK),
		Message: "All user successfully take it",
		Data:    respon,
	})
}

func (uh *UserHandler) DeleteUserById(c *echo.Context) error {
	id := c.Param("id")
	respon, err := uh.UserUseCase.DeleteUserById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, helper.DeleteUserTypeRespon{
		Status:  http.StatusText(http.StatusOK),
		Message: "User successfully deleted",
		Data:    respon,
	})
}

func (uh *UserHandler) UpdateUserById(c *echo.Context) error {
	var requestUser *user.UserUpdateModel
	err := c.Bind(&requestUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.Respon{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
	}
	id := c.Param("id")
	_, err = uh.UserUseCase.UpdateUserById(c.Request().Context(), id, *requestUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, helper.Respon{
		Status:  http.StatusText(http.StatusOK),
		Message: "User successfully updated",
	})
}
