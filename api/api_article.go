package api

import (
	"vgo/model"
	"vgo/service"

	"github.com/gin-gonic/gin"
)

// 创建文章
func CreateArticle(c *gin.Context) {
	cuser, _ := c.Get("user")
	if user, ok := cuser.(*model.User); ok {
		article := model.Article{}
		_ = c.ShouldBindJSON(&article)
		resp := article.Create(user)
		c.JSON(200, resp)
	}
}

// 删除文章
func DeleteArticle(c *gin.Context) {
	cuser, _ := c.Get("user")
	resp := &service.Response{}
	if user, ok := cuser.(*model.User); ok {
		article := model.Article{}
		_ = c.ShouldBindJSON(&article)
		resp = article.Delete(user)
	}
	c.JSON(200, resp)
}

// 获取文章详情
func GetArticleInfo(c *gin.Context) {
	id := c.Param("id")
	resp := model.GetArticleInfo(id)
	c.JSON(200, resp)
}

// 显示主页新闻
// 分别选取系部新闻和就业资讯的最后8个动态
func ShowIndex(c *gin.Context) {
	resp := model.GetArticleIndex()
	c.JSON(200, resp)
}

// 获取最新的新闻
func ShowArticleNews(c *gin.Context) {
	resp := model.GetArticleNews()
	c.JSON(200, resp)
}

// 获取文章列表
func GetArticleList(c *gin.Context) {
	getSubject := c.Param("subject")
	start := c.DefaultQuery("start", "0")
	limit := c.DefaultQuery("limit", "10")

	resp := model.GetArticleList(getSubject, start, limit)
	c.JSON(200, resp)
}

// 返回属于管用户的所有文章
// 管理员返回全部文章
func GetAdminArticleList(c *gin.Context) {
	user, _ := c.Get("user")
	if u, ok := user.(*model.User); ok {
		subject := c.DefaultQuery("subject", "")
		start := c.DefaultQuery("start", "0")
		limit := c.DefaultQuery("limit", "10")

		resp := model.AdminArticleList(u, subject, start, limit)
		c.JSON(200, resp)
	}
}

// 更新文章内容
func UpdateArticle(c *gin.Context) {
	type Info struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	user, _ := c.Get("user")
	if u, ok := user.(*model.User); ok {
		updateInfo := Info{}
		_ = c.BindJSON(&updateInfo)
		resp := model.UpdateArticle(u, updateInfo.ID, updateInfo.Title, updateInfo.Content)
		c.JSON(200, resp)
	}
}
