package repository

import (
	"context"
	// "fmt"
	"time"
	"ungraded-3/models"
	"ungraded-3/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageImplementation struct {
	DB *mongo.Collection
}

func NewMessageRepository(db *mongo.Collection) MessageImplementation {
	return MessageImplementation{
		DB: db,
	}
}

func (m MessageImplementation) Post(data *models.Message) *utils.ErrResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	res, err := m.DB.InsertOne(ctx, data)
	if err != nil {
		errResponse := &utils.ErrInternalServer
		errResponse.Detail = err.Error()
		return errResponse
	}
	data.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (m MessageImplementation) GetById(messageId string) *utils.ErrResponse {
	return nil
}

func (m MessageImplementation) GetBySenderReceiver(subject, email string) *utils.ErrResponse {
	return nil
}