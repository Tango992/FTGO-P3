package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Sender   string             `json:"sender" bson:"sender" validate:"required,email"`
	Receiver string             `json:"receiver" bson:"receiver" validate:"required,email"`
	Type     string             `json:"type" bson:"type" validate:"required"`
	Content  string             `json:"content" bson:"content" validate:"required"`
}
