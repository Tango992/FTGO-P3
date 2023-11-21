package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string             `bson:"name" json:"name"`
	Stock uint               `bson:"stock" json:"stock"`
}
