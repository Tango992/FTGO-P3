package repository

import (
	"context"
	"time"
	"ungraded-3/models"
	"ungraded-3/utils"

	"go.mongodb.org/mongo-driver/bson"
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

func (m MessageImplementation) GetById(messageId string) (models.Message,*utils.ErrResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()
	
	objectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		errResponse := &utils.ErrBadRequest
		errResponse.Detail = err.Error()
		return models.Message{}, errResponse
	}

	res := m.DB.FindOne(ctx, bson.M{"_id": objectId})

	var message models.Message
	if err := res.Decode(&message); err != nil {
		errResponse := &utils.ErrInternalServer
		errResponse.Detail = err.Error()
		return models.Message{}, errResponse
	}
	return message, nil
}

func (m MessageImplementation) GetBySenderReceiver(subject, email string) ([]models.Message, *utils.ErrResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	cursor, err := m.DB.Find(ctx, bson.M{subject: email})
	if err != nil {
		errResponse := &utils.ErrInternalServer
		errResponse.Detail = err.Error()
		return []models.Message{}, errResponse
	}
	
	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		errResponse := &utils.ErrInternalServer
		errResponse.Detail = err.Error()
		return []models.Message{}, errResponse
	}
	return messages, nil
}