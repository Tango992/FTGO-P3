package repository

import (
	"context"
	"preview-w2/dto"
	"preview-w2/models"
	"preview-w2/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return UserRepository{
		Collection: collection,
	}
}

func (u UserRepository) Register(data *models.User) *utils.ErrResponse {
	indexEmail := mongo.IndexModel{
		Keys: bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	
	if _, err := u.Collection.Indexes().CreateOne(context.TODO(), indexEmail); err != nil {
		return utils.ErrInternalServer.New(err.Error())
	}
	
	res, err := u.Collection.InsertOne(context.TODO(), data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.ErrConflict.New(err.Error())
		}
		return utils.ErrInternalServer.New(err.Error())
	}
	
	data.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (u UserRepository) FindUser(loginData dto.LoginUser) (models.User, *utils.ErrResponse) {
	result := u.Collection.FindOne(context.TODO(), bson.M{"email": loginData.Email})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, utils.ErrUnauthorized.New("Invalid email/password")
		}
		return models.User{}, utils.ErrInternalServer.New(err.Error())
	}
	
	var user models.User
	if err := result.Decode(&user); err != nil {
		return models.User{}, utils.ErrInternalServer.New(err.Error())
	}
	return user, nil
}