package main

import (
	"ungraded-2/config"
	"ungraded-2/helpers"
	"ungraded-2/routes"
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
