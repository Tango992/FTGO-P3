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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var employeeCollection = config.ConnectDB().Database("ungraded2DB").Collection("employees")

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
		Options: options.Index().SetUnique(true),
	}
	
	if _, err := employeeCollection.Indexes().CreateOne(ctx, indexEmail); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "Index Error", Data: &echo.Map{"error": err.Error()}})
	}
	
	if _, err := employeeCollection.InsertOne(ctx, user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "Insertion Error", Data: &echo.Map{"error": err.Error()}})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Registered",
		"data": user,
	})
}

func GetAllUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	page := c.QueryParam("page")
	if page != "" {
		return GetUsersPaginated(c, page)
	}
	
	res, err := employeeCollection.Find(ctx, echo.Map{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}
	defer res.Close(ctx)

	users := []models.User{}
	
	for res.Next(ctx) {
		var user models.User
		if err := res.Decode(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}
		users = append(users, user)
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

func GetUsersPaginated(c echo.Context, pageTmp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()
	
	page, err := strconv.Atoi(pageTmp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
	} else if page < 1 {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": "page must be greater than 0"}})
	}

	skip := 2 * (int64(page) - 1)
	opts := options.Find().SetLimit(2).SetSkip(skip)
	
	res, err := employeeCollection.Find(ctx, echo.Map{}, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}
	defer res.Close(ctx)

	users := []models.User{}
	
	for res.Next(ctx) {
		var user models.User
		if err := res.Decode(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
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
	
	res, err := employeeCollection.Find(ctx, echo.Map{}, options)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}
	defer res.Close(ctx)

	users := []models.User{}
	
	for res.Next(ctx) {
		var user models.User
		if err := res.Decode(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
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
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	var user models.User
	res := employeeCollection.FindOne(ctx, echo.Map{"_id": objectId})
	if res.Err() != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": res.Err().Error()}})
	}
	
	if err := res.Decode(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
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
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	if err := c.Validate(&user); err != nil {
		return err
	}

	userId := c.Param("userId")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
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

	res, err := employeeCollection.UpdateOne(ctx, bson.M{"_id": objectId}, updateData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

func DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 *time.Second)
	defer cancel()

	userId := c.Param("userId")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}
	
	res, err := employeeCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": res,
	})
}

