package routes

import (
	controller "preview-w2/controllers"
	"preview-w2/middlewares"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo, uc controller.UserController, pc controller.PaymentController) {
	e.POST("/users/register", uc.Register)
	e.POST("/users/login", uc.Login)

	e.POST("/payment", pc.CreatePayment, middlewares.RequireAuth)
}