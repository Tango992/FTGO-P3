package main

import (
	pb "grpc_test/internal/user"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	e := echo.New()
	e.POST("api/users", func(c echo.Context) error {
		req := new(pb.AddUserRequest)
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		res, err := client.AddUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"message": "Created",
			"data": res,
		})
	})

	e.GET("api/users", func(c echo.Context) error {
		res, err := client.GetUser(c.Request().Context(), &pb.GetUserRequest{})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Users",
			"data": res.Users,
		})
	})

	e.PUT("api/users/:id", func(c echo.Context) error {
		req := new(pb.UpdateUserRequest)
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		
		id := c.Param("id")
		req.Id = id
		
		res, err := client.UpdateUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Update users",
			"data": res,
		})
	})

	e.DELETE("api/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		req := &pb.DeleteUserRequest{Id: id}
		
		res, err := client.DeleteUser(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "User deleted successfully",
			"data": res,
		})
	})

	
    e.Logger.Fatal(e.Start(":8080"))
}