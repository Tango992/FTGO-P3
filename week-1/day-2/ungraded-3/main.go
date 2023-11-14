package main

import (
	"os"
	"ungraded-3/config"
	"ungraded-3/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	config.ConnectDB().Database("ungraded2DB")

	
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}