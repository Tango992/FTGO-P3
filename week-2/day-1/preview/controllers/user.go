package controllers

import (
	"net/http"
	"preview-w2/dto"
	"preview-w2/helpers"
	"preview-w2/models"
	"preview-w2/repository"
	"preview-w2/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	Repository repository.Users
}

func NewUserController(r repository.Users) UserController {
	return UserController{
		Repository: r,
	}
}

func (u UserController) Register(c echo.Context) error {
	var registerDataTmp dto.RegisterUser
	if err := c.Bind(&registerDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&registerDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := helpers.CreateHash(&registerDataTmp); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	registerData := models.User{
		Name: registerDataTmp.Name,
		Email: registerDataTmp.Email,
		Password: registerDataTmp.Password,
	}
	
	if err := u.Repository.Register(&registerData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	registerData.Password = ""
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Registered",
		Data: registerData,
	})
}

func (u UserController) Login(c echo.Context) error {
	var loginData dto.LoginUser
	if err := c.Bind(&loginData); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&loginData); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dbData, err := u.Repository.FindUser(loginData)
	if err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	if err := helpers.CheckPassword(dbData, loginData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	if err := helpers.SignNewJWT(c, dbData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Logged in",
		Data: "Authorization token is stored using cookie",
	})
}