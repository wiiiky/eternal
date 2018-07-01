package main

import (
	"context"
	"eternal/config"
	"eternal/controller"
	"eternal/errors"
	"eternal/event"
	"eternal/filemanager"
	"eternal/logging"
	"eternal/model/db"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const APPNAME = "eternal"

func main() {
	initConfig()
	initLogging()
	initDatabase()
	initEvent()
	filemanager.Init()
	initEcho(controller.Register)
}

func errorHandler(err error, c echo.Context) {
	if e, ok := err.(*errors.Error); ok {
		c.JSON(e.HttpStatus, e)
	} else if e, ok := err.(*echo.HTTPError); ok {
		c.JSON(e.Code, errors.NewError(0, -1, e.Message))
	} else {
		c.JSON(http.StatusInternalServerError, errors.NewError(0, -1, err.Error()))
	}
}

func initConfig() {
	config.Init(APPNAME)
}

/*
 * 初始化日志记录
 * https://github.com/sirupsen/logrus
 */
func initLogging() {
	format := config.GetStringDefault("log.format", "json")
	level := config.GetStringDefault("log.level", "info")
	output := config.GetStringDefault("log.output", "stdout")

	logging.Init(format, level, output)
}

/* 初始化数据库 */
func initDatabase() {
	dbURL := config.GetString("database.url")
	if dbURL == "" {
		log.Fatal("**CONFIG** database.url not found")
	} else if err := db.Init(dbURL); err != nil {
		log.Fatal("Connecting database failed:", err)
	}
}

/* 初始化事件发布模块 */
func initEvent() {
	amqpURL := config.GetString("event.amqp.url")
	if amqpURL == "" {
		log.Fatal("**CONFIG** event.amqp.url not found")
	}
	amqpExchange := config.GetString("event.amqp.exchange")
	if amqpExchange == "" {
		log.Fatal("**CONFIG** event.amqp.exchange not found")
	}
	amqpRouteKey := config.GetString("event.amqp.route_key")
	if amqpRouteKey == "" {
		log.Fatal("**CONFIG** event.amqp.route_key not found")
	}
	event.InitPub(amqpURL, amqpExchange, amqpRouteKey)
}

func initEcho(f func(*echo.Echo)) {
	httpAddr := config.GetString("http.addr")
	if httpAddr == "" {
		log.Fatal("**CONFIG** http.addr not found")
	}
	allowOrigins := config.GetStringSliceDefault("http.cors.origins", []string{"*"})
	allowMethods := config.GetStringSliceDefault("http.cors.methods", []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE})
	allowCredentials := config.GetBoolDefault("http.cors.credentials", false)
	sessionSecret := config.GetString("http.session.secret")
	accessLog := config.GetStringDefault("http.access_log", "stdout")
	if sessionSecret == "" {
		log.Fatal("**CONFIG** http.session.secret not found")
	}

	e := echo.New()
	accessLogWriter, err := logging.OpenLogFile(accessLog)
	if err != nil {
		log.Fatalf("**LOG** open %s failed: %s", accessLog, err)
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: accessLogWriter,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Source"},
		AllowCredentials: allowCredentials,
		AllowMethods:     allowMethods,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	e.HTTPErrorHandler = errorHandler

	f(e)

	go func() {
		if err := e.Start(httpAddr); err != nil {
			log.Info("shutting down the server:", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
