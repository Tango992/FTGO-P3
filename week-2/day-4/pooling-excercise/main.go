package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type User struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

var (
	dbPool *pgxpool.Pool
)

func dbConn() *pgxpool.Pool {
	dsn := "postgres://postgres:secret@localhost:5432/pooling?pool_max_conns=10"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

func main() {
	e := echo.New()

	e.POST("", createUser)
	e.Logger.Fatal(e.Start(":8080"))
}

func createUser(c echo.Context) error {
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer conn.Release()

	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	trx, err := conn.Begin(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := trx.QueryRow(context.Background(), "INSERT INTO USERS (name, age) VALUES $1, $2 RETURNING ID", user.Name, user.Age).Scan(user.Id); err != nil {
		trx.Rollback(context.Background())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User create",
		"user": user,
	})
}
