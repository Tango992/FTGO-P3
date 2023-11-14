package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name" validate:"required"`
	Address    string             `json:"address" bson:"address" validate:"required"`
	Email      string             `json:"email" bson:"email" validate:"required,email"`
	Salary     float32            `json:"salary" bson:"salary" validate:"required"`
	Department string             `json:"department" bson:"department" validate:"required"`
}
