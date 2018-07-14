package main

import (
	"context"
	"eternal/cache/store"
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
	initCache()
	initEvent()
	initFileManager()
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

func initFileManager() {
	filemanager.Init()
}

/* 初始化数据库 */
func initDatabase() {
	pgURL := config.GetString("database.pg.url")
	mongoURL := config.GetString("database.mongo.url")
	mongoDBName := config.GetString("database.mongo.dbname")
	if pgURL == "" {
		log.Fatal("**CONFIG** database.pg.url not found")
	} else if mongoURL == "" {
		log.Fatal("**CONFIG** database.mongo.url not found")
	} else if mongoDBName == "" {
		log.Fatal("**CONFIG** database.mongo.dbname not found")
	}
	if err := db.Init(pgURL, mongoURL, mongoDBName); err != nil {
		log.Fatal("Connecting database failed:", err)
	}
}

/* 初始化缓存 */
func initCache() {
	redisURL := config.GetString("cache.redis.url")
	if redisURL == "" {
		log.Fatal("**CONFIG** cache.redis.url not found")
	}
	if err := store.InitRedis(redisURL); err != nil {
		log.Fatal("InitRedis failed:", err)
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

/* 初始化HTTP服务 */
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

	/* 注册路由等 */
	f(e)

	go func() {
		if err := e.Start(httpAddr); err != nil {
			log.Info("shutting down the server:", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
