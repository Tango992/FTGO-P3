package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount    float32            `bson:"amount" json:"amount"`
	Channel   string             `bson:"channel" json:"channel"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt primitive.DateTime `bson:"completed_at" json:"completed_at"`
}
