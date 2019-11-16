package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"vgo/cache"
	"vgo/model"
	"vgo/service"
	"vgo/utils"

	"github.com/gin-gonic/gin"
)

// 获取验证码，有效期为 5 分钟
func HandleGetEmaiCode(c *gin.Context) {
	email := c.Query("email")
	resp := utils.SendCodeEmail(email)
	c.JSON(200, resp)
}

// 通过邮箱找回密码
/*
通过提交的数据获取验证码
判断缓存中折验证码是否过期
没过期，则提取其对应的邮箱，是否与提交的邮箱一致
一致则进入 DB，获取对应的用户，修改
如没有，则返回无此邮箱对应的用户
*/
type forgetPasswordForm struct {
	Email           string `json:"email"`
	Code            string `json:"code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func HandleForgetPassword(c *gin.Context) {
	form := forgetPasswordForm{}
	_ = c.BindJSON(&form)

	resp := &service.Response{}
	cacheEmail := cache.Get(form.Code)

	if cacheEmail == "" {
		resp.Code = 53004
		resp.Msg = "该验证码不存在或已过期"
	} else if cacheEmail != form.Email {
		resp.Code = 53003
		resp.Msg = "该邮箱与获取验证码的邮箱不一致"
	} else if form.Password != form.ConfirmPassword {
		resp.Code = 51002
		resp.Msg = "密码与确认密码不一致"
	} else {
		user := model.User{}
		model.DB.Where("email = ?", form.Email).First(&user)
		fmt.Println(user)
		if user.Id != 0 {
			resp = user.AdminUpdatePassword(form.Password)
		} else {
			resp.Code = 52002
			resp.Msg = "该邮箱还未注册"
		}
	}
	c.JSON(200, resp)

}

// 图片API
func HandleGetArticleImage(c *gin.Context) {
	filename := c.Param("filename")
	path := os.Getenv("IMAGE_PATH")

	// 返回类型
	contentType := ""
	if strings.LastIndex(filename, "png") != -1 {
		contentType = "image/png"
	} else {

		contentType = "image/jpeg"
	}

	if strings.Contains(filename, "/") == false {
		filepath := path + filename
		file, _ := os.Open(filepath)
		content, _ := ioutil.ReadAll(file)
		c.Data(200, contentType, content)
	}

}

// 获取用户
func getUser(c *gin.Context) *model.User {
	contextUser, _ := c.Get("user")
	user, _ := contextUser.(*model.User)
	return user
}
