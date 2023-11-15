package repository

import (
	"ungraded-3/models"

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

func (m MessageImplementation) Post(data *models.Message) error {
	return nil
}

func (m MessageImplementation) GetById(messageId string) error {
	return nil
}

func (m MessageImplementation) GetBySenderReceiver(subject, email string) error {
	return nil
}