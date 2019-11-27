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
	//engine.Use(middleware.Cors())
	// 大括号用于仅用于分隔代码，不作具体逻辑处理
	v1 := engine.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			go func() {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			}()
		})

		// 获取验证码
		v1.GET("/code", api.HandleGetEmaiCode)

		// 注册处理
		v1.POST("/register", api.HandleUserRegister)
		// 忘记密码
		v1.PUT("/password", api.HandleForgetPassword)

		// 用户登录处理
		v1.POST("/user", api.HandleUserLogin)

		v1.GET("/test", func(c *gin.Context) {
			s := sessions.Default(c)
			fmt.Println(s.Get("user_id"))
		})

		// 主页内容
		v1.GET("/article/index", api.HandleGetIndex)
		v1.GET("/article/info/:id", api.HandleGetArticleInfo)
		v1.GET("/article/newest", api.HandleGetArticleNews)

		// 图片API
		v1.GET("/img/:filename", api.HandleGetArticleImage)

		// bbs 列表
		v1.GET("/bbs/list", api.HandleGetBBSList)
		v1.GET("/bbs/info/:id", api.HandleGetBBSInfo)
		// 获取两个不同版块折文章列表
		v1.GET("/article/list/:subject", api.HandleGetArticleList)

		// 需要用户登录验证的 api
		auth := v1.Group("/")
		auth.Use(middleware.CurrentUser())
		auth.Use(middleware.AuthRequired())
		{
			// 用户相关
			auth.GET("user/logout", api.HandleUserLogout)
			auth.GET("user/info", api.HandleGetUserInfo)
			auth.POST("upload/avatar", api.HandleUploadAvatar)
			auth.GET("user/avatar/:filename", api.HandleGetAvatar)
			auth.GET("user/avatar", api.HandleGetAvatarNoFileName)

			// 微博 (bbs) 相关
			auth.POST("bbs/create", api.HandleCreateBBS)
			auth.DELETE("bbs/delete", api.HandleDeleteBBS)
			auth.PUT("bbs/add-comment", api.HandleAddComment)

			// 文章相关
			auth.POST("article/create", api.HandleCreateArticle)

			// 管理
			auth.GET("admin/article", api.HandleGetAdminArticleList)
			auth.GET("admin/bbs", api.HandleAdminBBSList)
			auth.PUT("admin/updateArticle", api.HandleUpdateArticle)
			auth.DELETE("admin/deleteArticle", api.HandleDeleteArticle)
			auth.PUT("/password/change", api.HandleLoginStatusChangePassword)
		}
	}

	return engine
}
