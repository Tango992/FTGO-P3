package main

import (
	"os"
	"ungraded-2/config"
	"ungraded-2/controllers"
	"ungraded-2/helpers"
	"ungraded-2/repository"
	"ungraded-2/routes"

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
	
	db := config.ConnectDB().Database("ungraded2DB")
	employeeCollection := db.Collection("employees")
	
	employeeRepository := repository.NewEmployeeRepository(employeeCollection)
	employeeController := controllers.NewEmployeeController(employeeRepository)
	
	routes.UserRoute(e, employeeController)
	
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
