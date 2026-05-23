package domain

import (
	"context"
	"p3-lc01-billyhaffas/internal/model/user"

	"github.com/labstack/echo/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	PostUser(ctx context.Context, user *user.UserInsertModel) error
	GetAllUserEmail(ctx context.Context, email string) (*user.UserEmail, error)
	GetUserById(ctx context.Context, userId primitive.ObjectID) (*user.UserSelectModel, error)
	GetAllUser(ctx context.Context) ([]user.UserSelectModel, error)
	DeleteUserById(ctx context.Context, userId primitive.ObjectID) (*mongo.DeleteResult, error)
	UpdateUserById(ctx context.Context, userId primitive.ObjectID, updateRequest user.UserUpdateModel) (*mongo.UpdateResult, error)
}

type UserUseCase interface {
	PostUser(ctx context.Context, request *user.UserInsertModel) error
	GetUserById(ctx context.Context, userId string) (*user.UserSelectModel, error)
	GetAllUser(ctx context.Context) ([]user.UserSelectModel, error)
	DeleteUserById(ctx context.Context, userId string) (*mongo.DeleteResult, error)
	UpdateUserById(ctx context.Context, userId string, updateRequest user.UserUpdateModel) (*mongo.UpdateResult, error)
}

type UserHandler interface {
	PostUser(c *echo.Context) error
	GetUserById(c *echo.Context) error
	GetAllUser(c *echo.Context) error
	DeleteUserById(c *echo.Context) error
	UpdateUserById(c *echo.Context) error
}
