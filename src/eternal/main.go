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
	initEcho(func(e *echo.Echo) {
		e.PUT("/account", view.Login)
		e.POST("/account", view.Signup)
		e.GET("/supported_countries", view.GetSupportedCountries)

		g := e.Group("", cmiddleware.AuthMiddleware)
		g.GET("/account", view.GetAccountInfo)
		g.GET("/user/profile", view.GetUserProfile)
	})
}

func errorHandler(err error, c echo.Context) {
	if e, ok := err.(*view.Error); ok {
		c.JSON(e.HttpStatus, e)
	} else if e, ok := err.(*echo.HTTPError); ok {
		c.JSON(e.Code, view.NewError(0, -1, e.Message))
	} else {
		c.JSON(http.StatusInternalServerError, view.NewError(0, -1, err.Error()))
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

func initEcho(f func(*echo.Echo)) {
	viper.SetDefault("http.cors.methods", []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE})
	viper.SetDefault("http.cors.origins", []string{"*"})
	viper.SetDefault("http.cors.credentials", false)

	httpAddr := viper.GetString("http.addr")
	allowOrigins := viper.GetStringSlice("http.cors.origins")
	allowMethods := viper.GetStringSlice("http.cors.methods")
	allowCredentials := viper.GetBool("http.cors.credentials")
	sessionSecret := viper.GetString("http.session.secret")
	log.Debugf("http.addr:%s", httpAddr)
	log.Debugf("http.cors.origins:%s", allowOrigins)
	log.Debugf("http.cors.methods:%s", allowMethods)
	log.Debugf("http.cors.credentials:%v", allowCredentials)
	log.Debugf("http.session.secret:%s", sessionSecret)
	if httpAddr == "" {
		log.Fatal("Incomplete config. http.addr not found")
	}
	if len(allowOrigins) == 0 {
		log.Fatal("Incomplete config. http.cors.origins not found")
	}
	if len(allowMethods) == 0 {
		log.Fatal("Incomplete config. http.cors.methods not found")
	}
	if sessionSecret == "" {
		log.Fatal("Incomplete config. http.session.secret not found")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: allowCredentials,
		AllowMethods:     allowMethods,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	e.HTTPErrorHandler = errorHandler

	f(e)

	e.Logger.Fatal(e.Start(httpAddr))
}
