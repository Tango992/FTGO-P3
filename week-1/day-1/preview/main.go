package main

import (
	"log"
	"os"
	"preview-w1/config"
	"preview-w1/controller"
	"preview-w1/entity"
	"preview-w1/helpers"
	"preview-w1/middlewares"
	"preview-w1/repository"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func main() {
	db := config.InitDb()
	if err := db.AutoMigrate(entity.User{}, entity.Transaction{}, entity.Loan{}); err != nil {
		log.Fatal(err)
	}
	
	dbHandler := repository.NewDbHandler(db)
	userController := controller.NewUserController(dbHandler)
	loanController := controller.NewLoanController(dbHandler)
	
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	logger := zerolog.New(os.Stdout)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Str("Method", v.Method).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))	
	e.Use(middleware.Recover())
	
	v1 := e.Group("/v1/ms-paylater")
	{
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)
		requireAuth := v1.Group("")
		requireAuth.Use(middlewares.RequireAuth)
		{
			requireAuth.POST("/loan", loanController.ProposeLoan)
			requireAuth.GET("/limit", loanController.Limit)
			requireAuth.POST("/tarik-saldo", loanController.Withdraw)
			requireAuth.POST("/pay", loanController.Deposit)
		}
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
