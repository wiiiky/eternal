package main

import (
	"context"
	"eternal/errors"
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
	"github.com/spf13/viper"
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
	initConfig()
	initLogging()
	initDatabase()
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
		authApi.GET("/question/:id", questionView.GetQuestion)

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
 * 读取配置
 * https://github.com/spf13/viper
 */
func initConfig() {
	viper.SetConfigName("eternal")           // name of config file (without extension)
	viper.AddConfigPath("/etc/" + APPNAME)   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + APPNAME) // call multiple times to add many search paths
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(err)
	}
	viper.SetDefault("debug", true)
}

/*
 * 初始化日志记录
 * https://github.com/sirupsen/logrus
 */
func initLogging() {
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "stdout")

	logging.Init(viper.GetString("log.format"), viper.GetString("log.level"), viper.GetString("log.output"))
}

/* 初始化数据库 */
func initDatabase() {
	dbURL := viper.GetString("database.url")
	if err := db.Init(dbURL); err != nil {
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
