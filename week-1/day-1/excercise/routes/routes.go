package routes

import (
	"excercise-mongo/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/user", controllers.CreateUser)
	e.GET("/user:userId", controllers.GetUserById)
	e.GET("/user", controllers.GetAllUser)
	e.PUT("/user:userId", controllers.UpdateUser)
	e.DELETE("/user", controllers.DeleteUser)
}
