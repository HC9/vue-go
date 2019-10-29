package route

import (
	"fmt"
	"vgo/api"
	"vgo/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()
	// 添加中间件
	engine.Use(middleware.Session())
	engine.Use(middleware.Cors())
	// 大括号用于仅用于分隔代码，不作具体逻辑处理
	v1 := engine.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// 获取验证码
		v1.GET("/code", api.GetCode)

		// 注册处理
		v1.POST("/user/register", api.UserRegister)
		// 用户邮件验证处理
		v1.GET("/user/mail/:key", api.UserInsert)
		// 忘记密码
		v1.PUT("/forget-password", api.HandleForgetPassword)

		// 用户登录处理
		v1.POST("/user/login", api.UserLogin)

		v1.GET("/test", func(c *gin.Context) {
			s := sessions.Default(c)
			fmt.Println(s.Get("user_id"))
		})

		// 主页内容
		v1.GET("/article/index", api.ShowIndex)
		v1.GET("/article/info/:id", api.GetArticleInfo)
		v1.GET("/article/newest", api.ShowArticleNews)

		// 图片API
		v1.GET("/img/:filename", api.HandleGetImage)

		// bbs 列表
		v1.GET("/bbs/list", api.HandleBBSList)
		v1.GET("/bbs/info/:id", api.HandleBBSInfo)
		// 获取两个不同版块折文章列表
		v1.GET("/article/list/:subject", api.GetArticleList)

		// 需要用户登录验证的 api
		auth := v1.Group("/")
		auth.Use(middleware.CurrentUser())
		auth.Use(middleware.AuthRequired())
		{
			// 用户相关
			auth.GET("user/logout", api.UserLogout)
			auth.GET("user/info", api.GetUserInfo)
			auth.POST("upload/avatar", api.HandleUploadAvatar)
			auth.GET("user/avatar/:filename", api.HandleGetAvatar)
			auth.GET("user/avatar", api.HandleGetAvatarNoFileName)

			// 微博 (bbs) 相关
			auth.POST("bbs/create", api.CreateBBS)
			auth.DELETE("bbs/delete", api.DeleteBBS)
			auth.PUT("bbs/add-comment", api.AddComment)

			// 文章相关
			auth.POST("article/create", api.CreateArticle)

			// 管理
			auth.GET("admin/article", api.GetAdminArticleList)
			auth.GET("admin/bbs", api.HandleAdminBBSList)
			auth.PUT("admin/updateArticle", api.UpdateArticle)
			auth.DELETE("admin/deleteArticle", api.DeleteArticle)
			auth.PUT("admin/change-password", api.AdminChangePassword)
		}
	}

	return engine
}
