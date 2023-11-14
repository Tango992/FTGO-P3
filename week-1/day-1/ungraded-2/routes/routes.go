package routes

import (
	"ungraded-2/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo, employeeController controllers.EmployeeController) {
	e.POST("/employees", employeeController.CreateUser)
	e.GET("/employees/:userId", employeeController.GetUserById)
	e.GET("/employees", employeeController.GetAllUser)
	e.GET("/employees/salary", employeeController.GetAllUserSortedBySalary)
	e.PUT("/employees/:userId", employeeController.UpdateUser)
	e.DELETE("/employees/:userId", employeeController.DeleteUser)
}
