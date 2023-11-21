package main

import (
	"net/http"
	"ungraded_5/config"
	"ungraded_5/controller"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func main() {
	conn, client := config.InitGrpcClient()
	defer conn.Close()

	productController := controller.NewProductController(client)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/products", productController.Create)
	e.GET("/products", productController.GetAll)
	e.PUT("/products/:id", productController.Update)
	e.DELETE("/products/:id", productController.Delete)
	
	e.Logger.Fatal(e.Start(":8080"))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
  }