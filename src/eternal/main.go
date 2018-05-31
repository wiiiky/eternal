package main

import (
	"context"
	"eternal/config"
	"eternal/errors"
	"eternal/event"
	"eternal/filemanager"
	"eternal/logging"
	cmiddleware "eternal/middleware"
	"eternal/model/db"
	accountView "eternal/view/account"
	fileView "eternal/view/file"
	homeView "eternal/view/home"
	questionView "eternal/view/question"
	userView "eternal/view/user"
	"github.com/go-playground/validator"
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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config.Init(APPNAME)
	initLogging()
	initDatabase()
	initEvent()
	filemanager.Init()
	initEcho(func(e *echo.Echo) {
		e.Validator = &CustomValidator{validator: validator.New()}

		api := e.Group("/api")

		// 登录注册
		api.PUT("/account/token", accountView.Login)                       // 登录
		api.POST("/account", accountView.Signup)                           // 注册
		api.GET("/supported_countries", accountView.GetSupportedCountries) // 获取支持的国家列表
		api.GET("/file/:id", fileView.DownloadFile)                        // 下载文件

		authApi := api.Group("", cmiddleware.AuthMiddleware)
		authApi.DELETE("/account/token", accountView.Logout) // 注销
		authApi.GET("/account", accountView.GetAccountInfo)  // 获取账号信息
		// 用户相关
		authApi.GET("/user/profile", userView.GetUserProfile) // 获取用户信息
		authApi.PUT("/user/cover", userView.UpdateUserCover)  // 更新用户的封面图
		// 主页相关
		authApi.GET("/home/hot/answers", homeView.GetHotAnswers) // 获取热门回答
		// 回答相关
		authApi.POST("/answer/:id/upvote", questionView.UpvoteAnswer)
		authApi.POST("/answer/:id/downvote", questionView.DownvoteAnswer)
		authApi.DELETE("/answer/:id/upvote", questionView.UndoUpvoteAnswer)
		authApi.DELETE("/answer/:id/downvote", questionView.UndoDownvoteAnswer)
		// 话题相关
		authApi.GET("/topics", questionView.FindTopics)
		// 问题相关
		authApi.POST("/question", questionView.CreateQuestion)
		authApi.GET("/questions", questionView.FindQuestions)
		authApi.GET("/question/:id", questionView.GetQuestion)
		authApi.GET("/question/:qid/answers", questionView.GetQuestionAnswers)

		// 上传文件
		authApi.POST("/file", fileView.UploadFile)
	})
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
	if sessionSecret == "" {
		log.Fatal("**CONFIG** http.session.secret not found")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: allowCredentials,
		AllowMethods:     allowMethods,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	e.HTTPErrorHandler = errorHandler

	f(e)

	go func() {
		if err := e.Start(httpAddr); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
