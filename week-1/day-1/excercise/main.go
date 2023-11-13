package main

import (
	"excercise-mongo/config"
	"excercise-mongo/helpers"
	"excercise-mongo/routes"
	"os"
	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.ConnectDB()
	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.UserRoute(e)
	
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
