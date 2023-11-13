package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"ungraded-2/config"
	"ungraded-2/helpers"
	"ungraded-2/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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
	
	res, err := UserCollection.InsertOne(ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Registered",
		"data": res,
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

func GetUserById(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	userIdTmp := c.Param("userId")
	userId, err := strconv.Atoi(userIdTmp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	var user models.User
	res := UserCollection.FindOne(ctx, echo.Map{"_id": userId})
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

	userIdTmp := c.Param("userId")
	userId, err := strconv.Atoi(userIdTmp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
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

	userIdTmp := c.Param("userId")
	userId, err := strconv.Atoi(userIdTmp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	res, err := UserCollection.DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

