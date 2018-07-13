package controller

import (
	"eternal/controller/account"
	"eternal/controller/file"
	"eternal/controller/middleware"
	"eternal/controller/misc"
	"eternal/controller/question"
	"eternal/controller/user"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Register(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}

	api := e.Group("/api", middleware.SourceMiddleware)

	// 登录注册
	api.PUT("/account/token", account.Login)                       // 登录
	api.POST("/account", account.Signup)                           // 注册
	api.GET("/supported_countries", account.GetSupportedCountries) // 获取支持的国家列表
	api.GET("/file/:id", file.DownloadFile)                        // 下载文件
	api.POST("/sms/signup", misc.SendSignupCode)                   // 发送注册短信

	authApi := api.Group("", middleware.AuthMiddleware)
	authApi.DELETE("/account/token", account.Logout) // 注销
	authApi.GET("/account", account.GetAccountInfo)  // 获取账号信息
	// 用户相关
	authApi.GET("/user/profile", user.GetUserProfile) // 获取用户信息
	authApi.PUT("/user/cover", user.UpdateUserCover)  // 更新用户的封面图
	// 回答相关
	authApi.GET("/hot/answers", question.GetHotAnswers) // 获取热门回答
	authApi.POST("/answer/:id/upvote", question.UpvoteAnswer)
	authApi.POST("/answer/:id/downvote", question.DownvoteAnswer)
	authApi.DELETE("/answer/:id/upvote", question.UndoUpvoteAnswer)
	authApi.DELETE("/answer/:id/downvote", question.UndoDownvoteAnswer)
	// 话题相关
	authApi.GET("/topics", question.FindTopics)
	// 问题相关
	authApi.POST("/question", question.CreateQuestion)
	authApi.GET("/questions", question.FindQuestions)
	authApi.GET("/question/:id", question.GetQuestion)
	authApi.GET("/question/:qid/answers", question.GetQuestionAnswers)
	authApi.POST("/question/:id/follow", question.FollowQuestion)
	authApi.DELETE("/question/:id/follow", question.UnfollowQuestion)

	// 上传文件
	authApi.POST("/file", file.UploadFile)
}
