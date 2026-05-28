package usecase

import (
	"context"
	"errors"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUseCase struct {
	UserRepository domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &UserUseCase{UserRepository: userRepo}
}

func (uc *UserUseCase) PostUser(ctx context.Context, request *user.UserInsertModel) error {
	if request.Email == "" {
		return errors.New("email can't empty")
	}
	existingEmail, err := uc.UserRepository.GetAllUserEmail(ctx, request.Email)
	if existingEmail != nil {
		return errors.New("email already exists")
	}

	if err != nil {
		return err
	}

	err = uc.UserRepository.PostUser(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) GetUserById(ctx context.Context, userId string) (*user.UserSelectModel, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	userRespon, err := uc.UserRepository.GetUserById(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return userRespon, nil
}

func (uc *UserUseCase) GetAllUser(ctx context.Context) ([]user.UserSelectModel, error) {
	userRespon, err := uc.UserRepository.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}
	return userRespon, nil
}

func (uc *UserUseCase) DeleteUserById(ctx context.Context, userId string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	userRespon, err := uc.UserRepository.GetUserById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if userRespon == nil {
		return nil, errors.New("user doesn't exist on database")
	}
	deleteRespon, err := uc.UserRepository.DeleteUserById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return deleteRespon, nil
}

func (uc *UserUseCase) UpdateUserById(ctx context.Context, userId string, updateRequest user.UserUpdateModel) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	userRespon, err := uc.UserRepository.GetUserById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if userRespon == nil {
		return nil, errors.New("user doesn't exist on database")
	}
	updateRespon, err := uc.UserRepository.UpdateUserById(ctx, objectID, updateRequest)
	if err != nil {
		return nil, err
	}
	return updateRespon, nil
}
