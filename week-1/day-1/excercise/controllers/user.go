package controllers

import (
	"context"
	"excercise-mongo/config"
	"excercise-mongo/helpers"
	"excercise-mongo/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserCollection = config.ConnectDB().Database("db_rmt02").Collection("users")

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
	
	newUser := models.User{
		Id: primitive.NewObjectID(),
		Name: user.Name,
		Email: user.Email,
		Address: user.Address,
	}

	_, err := UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Registered",
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
	
	return c.JSON(http.StatusCreated, echo.Map{
		"data": users,
	})
}

func GetUserById(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	userId := c.Param("userId")
	// objId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})

	// }
	var user models.User
	res := UserCollection.FindOne(ctx, echo.Map{"_id": userId})
	if res.Err() != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": res.Err().Error()}})
	}
	
	if err := res.Decode(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, echo.Map{
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

	updateData := bson.M{
		"$set": bson.M{
			"name": user.Name,
			"address": user.Address,
			"email": user.Email,
		},
	}

	res, err := UserCollection.UpdateOne(ctx, bson.M{"_id": userId}, updateData)
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

	res, err := UserCollection.DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

