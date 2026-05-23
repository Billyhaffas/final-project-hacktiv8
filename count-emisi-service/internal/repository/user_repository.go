package repository

import (
	"context"
	"fmt"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserCollection(collectionRepo *mongo.Collection) domain.UserRepository {
	return &UserRepository{Collection: collectionRepo}
}

func (cp *UserRepository) PostUser(ctx context.Context, user *user.UserInsertModel) error {
	_, err := cp.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (cp *UserRepository) GetAllUserEmail(ctx context.Context, email string) (*user.UserEmail, error) {
	var emailData user.UserEmail
	err := cp.Collection.FindOne(ctx,
		bson.M{
			"email": email,
		},
	).Decode(&emailData)
	fmt.Println(emailData)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &emailData, nil

}

func (cp *UserRepository) GetUserById(ctx context.Context, userId primitive.ObjectID) (*user.UserSelectModel, error) {
	var userData user.UserSelectModel
	err := cp.Collection.FindOne(ctx, bson.M{
		"_id": userId,
	}).Decode(&userData)
	fmt.Println(userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err

	}
	return &userData, nil
}

func (cp *UserRepository) GetAllUser(ctx context.Context) ([]user.UserSelectModel, error) {
	var usersData []user.UserSelectModel
	cursor, err := cp.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &usersData)
	fmt.Println(usersData)
	fmt.Println("Collection:", cp.Collection.Name())
	if err != nil {
		return nil, err
	}
	return usersData, nil
}

func (cp *UserRepository) DeleteUserById(ctx context.Context, userId primitive.ObjectID) (*mongo.DeleteResult, error) {
	respon, err := cp.Collection.DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		return nil, err
	}
	return respon, nil
}

func (cp *UserRepository) UpdateUserById(ctx context.Context, userId primitive.ObjectID, updateRequest user.UserUpdateModel) (*mongo.UpdateResult, error) {
	respon, err := cp.Collection.UpdateByID(ctx, userId, bson.M{"$set": updateRequest})
	if err != nil {
		return nil, err
	}
	return respon, nil
}
