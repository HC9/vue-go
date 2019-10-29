package middleware

import (
	"strconv"
	"vgo/model"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		if id := s.Get("user_id"); id != nil {
			convertId := strconv.Itoa(id.(int))
			user := model.GetOneUser(convertId)
			c.Set("user", user)
			c.Next()
		}
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user"); user != nil {
			// 断言，判断其是否是 model 数据类型
			if _, ok := user.(*model.User); ok {
				//fmt.Println(user)
				c.Next()
				return
			}
		}
		// 不是用户类型，则返回，Abort 中断执行流程返回
		c.JSON(200, gin.H{
			"status": 50001,
			"msg":    "该操作需要先登录",
		})
		c.Abort()
	}
}
