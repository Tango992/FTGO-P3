package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"ungraded-2/dto"
	"ungraded-2/helpers"
	"ungraded-2/models"
	"ungraded-2/repository"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmployeeController struct {
	Repository repository.Employee
}

func NewEmployeeController(r repository.Employee) EmployeeController {
	return EmployeeController{
		Repository: r,
	}
}

func (e EmployeeController) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if err := c.Validate(&user); err != nil {
		return err
	}

	if err := e.Repository.Create(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Registered",
		"data": user,
	})
}

func (e EmployeeController) GetAllUser(c echo.Context) error {
	page := c.QueryParam("page")
	if page != "" {
		return e.GetUsersPaginated(c, page)
	}
	
	users, err := e.Repository.Find(nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

func (e EmployeeController) GetUsersPaginated(c echo.Context, pageTmp string) error {
	page, err := strconv.Atoi(pageTmp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
	} else if page < 1 {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": "page must be greater than 0"}})
	}

	skip := 2 * (int64(page) - 1)
	opts := options.Find().SetLimit(2).SetSkip(skip)
	
	users, err := e.Repository.Find(opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

func (e EmployeeController) GetAllUserSortedBySalary(c echo.Context) error {
	opts := options.Find().SetSort(bson.D{{Key: "salary", Value: 1}})
	
	users, err := e.Repository.Find(opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
	})
}

func (e EmployeeController) GetUserById(c echo.Context) error {
	userId := c.Param("userId")

	user, err := e.Repository.FindById(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": user,
	})
}

func (e EmployeeController) UpdateUser(c echo.Context) error {
	var user dto.UpdateEmployee
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	if err := c.Validate(&user); err != nil {
		return err
	}

	userId := c.Param("userId")
	if err := e.Repository.Update(user, userId); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"data": "Updated successfully",
	})
}

func (e EmployeeController) DeleteUser(c echo.Context) error {
	userId := c.Param("userId")
	if err := e.Repository.DeleteById(userId); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": fmt.Sprintf("_id: %v deleted successfully", userId),
	})
}

