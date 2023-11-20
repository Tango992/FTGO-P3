package repository

import (
	"context"
	"preview-w2/models"
	"preview-w2/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	Collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) PaymentRepository {
	return PaymentRepository{
		Collection: collection,
	}
}

func (p PaymentRepository) Create(data *models.Payment) *utils.ErrResponse {
	res, err := p.Collection.InsertOne(context.TODO(), data)
	if err != nil {
		return utils.ErrInternalServer.New(err.Error())
	}
	
	data.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}