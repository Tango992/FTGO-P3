package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" validate:"required"`
	Address string             `json:"address" validate:"required"`
	Email   string             `json:"email" validate:"required"`
}
