package routes

import (
	"ungraded-2/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/employees", controllers.CreateUser)
	e.GET("/employees/:userId", controllers.GetUserById)
	e.GET("/employees", controllers.GetAllUser)
	e.GET("/employees/salary", controllers.GetAllUserSortedBySalary)
	e.PUT("/employees/:userId", controllers.UpdateUser)
	e.DELETE("/employees/:userId", controllers.DeleteUser)
}
