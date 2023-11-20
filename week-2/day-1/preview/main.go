package main

import (
	"os"
	"preview-w2/config"
	"preview-w2/controllers"
	"preview-w2/helpers"
	"preview-w2/repository"
	"preview-w2/routes"
	"preview-w2/scheduler"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := config.ConnectDB().Database("preview_w2_db")
	
	usersCollection := db.Collection("users")
	usersRepository := repository.NewUserRepository(usersCollection)
	usersController := controllers.NewUserController(usersRepository)

	paymentsCollection := db.Collection("payments")
	paymentsRepository := repository.NewPaymentRepository(paymentsCollection)
	paymentsController := controllers.NewPaymentController(paymentsRepository)

	routes.Routes(e, usersController, paymentsController)

	cron := scheduler.NewScheduler(db)
	cron.StartCronJob()

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}