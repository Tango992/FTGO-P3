package repository

import (
	"context"
	"time"
	"ungraded-2/dto"
	"ungraded-2/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmployeeImplementation struct {
	DB *mongo.Collection
}

func NewEmployeeRepository(db *mongo.Collection) EmployeeImplementation {
	return EmployeeImplementation{
		DB: db,
	}
}

func (p EmployeeImplementation) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()
	
	indexEmail := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	
	if _, err := p.DB.Indexes().CreateOne(ctx, indexEmail); err != nil {
		return err
	}
	
	if _, err := p.DB.InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}

func (p EmployeeImplementation) Find(opts *options.FindOptions) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()
	
	res, err := p.DB.Find(ctx, bson.M{}, opts)
	if err != nil {
		return []models.User{}, err
	}
	defer res.Close(ctx)

	users := []models.User{}
	
	for res.Next(ctx) {
		var user models.User
		if err := res.Decode(&user); err != nil {
			return []models.User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (p EmployeeImplementation) FindById(userId string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	res := p.DB.FindOne(ctx, bson.M{"_id": objectId})
	if res.Err() != nil {
		return models.User{}, err
	}
	
	if err := res.Decode(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (p EmployeeImplementation) Update(user dto.UpdateEmployee, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	updateData := bson.M{
		"$set": user,
	}

	_, err = p.DB.UpdateOne(ctx, bson.M{"_id": objectId}, updateData)
	if err != nil {
		return err
	}
	return nil
}

func (p EmployeeImplementation) DeleteById(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	_, err = p.DB.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}
	
	return nil
}