package main

import (
	"eternal/logging"
	cmiddleware "eternal/middleware"
	"eternal/model/db"
	"eternal/view"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

const APPNAME = "eternal"

func main() {
	initConfig()
	initLogging()
	initDatabase()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:1323"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("xGKCAxcCbUZyM3GayitHcQJz9HHnDNKk"))))
	e.Use(cmiddleware.AuthMiddleware)
	e.HTTPErrorHandler = errorHandler

	e.PUT("/login", view.Login)
	e.PUT("/signup", view.Signup)
	e.GET("/user", view.GetAccountInfo)

	e.Logger.Fatal(e.Start(":1323"))
}

func errorHandler(err error, c echo.Context) {
	if e, ok := err.(*view.Error); ok {
		c.JSON(e.HttpStatus, e)
	} else if e, ok := err.(*echo.HTTPError); ok {
		c.JSON(e.Code, view.NewError(0, -1, e.Message))
	} else {
		c.JSON(http.StatusInternalServerError, view.NewError(0, -1, e.Error()))
	}
}

/*
 * 读取配置
 * https://github.com/spf13/viper
 */
func initConfig() {
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.AddConfigPath("/etc/" + APPNAME)   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + APPNAME) // call multiple times to add many search paths
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(err)
	}
}

/*
 * 初始化日志记录
 * https://github.com/sirupsen/logrus
 */
func initLogging() {
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "stdout")

	logging.Start(viper.GetString("log.format"), viper.GetString("log.level"), viper.GetString("log.output"))
}

/* 初始化数据库 */
func initDatabase() {
	dbURL := viper.GetString("database.url")
	if err := db.Start(dbURL); err != nil {
		log.Fatal("Connecting database failed:", err)
	}
}
