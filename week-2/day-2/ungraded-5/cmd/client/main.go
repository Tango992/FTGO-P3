package main

import (
	"excercise-mongo/helpers"
	"ungraded_5/config"
	"ungraded_5/controller"
	"ungraded_5/service"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {
	mbConn, ch := config.InitMessageBroker()
	defer ch.Close()
	defer mbConn.Close()
	
	conn, grpcClient := config.InitGrpcClient()
	defer conn.Close()

	redisClient := config.InitCache()

	messageBrokerService := service.NewMessageBrokerService(ch)
	productController := controller.NewProductController(grpcClient, redisClient, messageBrokerService)

	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/products", productController.Create)
	e.GET("/products", productController.GetAll)
	e.PUT("/products/:id", productController.Update)
	e.DELETE("/products/:id", productController.Delete)

	e.Logger.Fatal(e.Start(":8080"))
}