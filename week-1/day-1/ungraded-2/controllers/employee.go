package controllers

import (
	"context"
	"net/http"
	"time"
	"ungraded-2/config"
	"ungraded-2/helpers"
	"ungraded-2/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection = config.ConnectDB().Database("ungraded2DB").Collection("employees")

func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if err := c.Validate(&user); err != nil {
		return err
	}

	indexEmail := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
	}
	
	if _, err := UserCollection.Indexes().CreateOne(ctx, indexEmail); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	if _, err := UserCollection.InsertOne(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Registered",
		"data": user,
	})
}

func GetAllUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	res, err := UserCollection.Find(ctx, echo.Map{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	defer res.Close(ctx)

	users := []models.User{}
	
	for res.Next(ctx) {
		var user models.User
		if err := res.Decode(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
		users = append(users, user)
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

func GetAllUserSortedBySalary(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()
	options := options.Find().SetSort(bson.D{{Key: "salary", Value: 1}})
	
	res, err := UserCollection.Find(ctx, echo.Map{}, options)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	defer res.Close(ctx)

	users := []models.User{}
	
	for res.Next(ctx) {
		var user models.User
		if err := res.Decode(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
		users = append(users, user)
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

func GetUserById(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	userId := c.Param("userId")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var user models.User
	res := UserCollection.FindOne(ctx, echo.Map{"_id": objectId})
	if res.Err() != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": res.Err().Error()}})
	}
	
	if err := res.Decode(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": user,
	})
}

func UpdateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if err := c.Validate(&user); err != nil {
		return err
	}

	userId := c.Param("userId")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updateData := bson.M{
		"$set": bson.M{
			"name": user.Name,
			"address": user.Address,
			"email": user.Email,
			"salary": user.Salary,
			"department": user.Department,
		},
	}

	res, err := UserCollection.UpdateOne(ctx, bson.M{"_id": objectId}, updateData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

func DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if err := c.Validate(&user); err != nil {
		return err
	}

	userId := c.Param("userId")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	res, err := UserCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

