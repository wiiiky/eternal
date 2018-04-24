package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"eternal/db"
	"fmt"
)

func main() {
	Init()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func Init() {
	err := db.Start("postgres://postgres@127.0.0.1/test")
	fmt.Println(err)
}
