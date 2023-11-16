package main

import (
	"os"
	"ungraded-3/config"
	"ungraded-3/controller"
	"ungraded-3/helpers"
	"ungraded-3/repository"
	"ungraded-3/routes"
	"ungraded-3/scheduler"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	db := config.ConnectDB().Database("ungraded3DB")
	messageCollection := db.Collection("messages")

	messageRepository := repository.NewMessageRepository(messageCollection)
	messageController := controller.NewMessageController(messageRepository)
	routes.UserRoute(e, messageController)

	schedule := scheduler.New(messageCollection)
	schedule.StarWorker()
	
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}