package api

import (
	"vgo/model"
	"vgo/service"

	"github.com/gin-gonic/gin"
)

// 创建微博
func HandleCreateBBS(c *gin.Context) {
	user := getUser(c)
	bbs := model.BBS{}
	bbs.Avatar = "/api/v1/user/avatar/" + user.Avatar
	_ = c.ShouldBindJSON(&bbs)
	bbs.Username = user.Username
	resp := bbs.Create()
	c.JSON(200, resp)
}

// 添加评论
func HandleAddComment(c *gin.Context) {
	user := getUser(c)
	postComment := &model.PostComment{}
	postComment.Username = user.Username
	_ = c.BindJSON(postComment)
	postComment.AddComment(user.Avatar)

	c.JSON(200, gin.H{
		"code": 20000,
		"msg":  "评论添加成功",
	})
}

// 删除微博
func HandleDeleteBBS(c *gin.Context) {
	user := getUser(c)

	w := &model.BBS{}
	_ = c.ShouldBindJSON(w)
	// 0 号权限是管理员，能够不用验证执行删帖操作
	resp := &service.Response{}
	if user.Role == 0 {
		resp = w.Delete()
	} else {
		// 对用户名进行验证
		info := w.GetBBSInfo(w.ID)
		if temp, ok := info.Data.(*model.BBS); ok {
			if temp.Username == user.Username {
				resp = w.Delete()
			} else {
				resp.Code = 55003
				resp.Msg = "你没有权限执行此操作"
			}
		}
	}

	c.JSON(200, resp)
}

// 微博详情
func HandleGetBBSInfo(c *gin.Context) {
	id := c.Param("id")
	w := &model.BBS{}
	resp := w.GetBBSInfo(id)
	c.JSON(200, resp)
}

// 列表
func HandleGetBBSList(c *gin.Context) {
	start := c.DefaultQuery("start", "0")
	limit := c.DefaultQuery("limit", "10")
	resp := model.GetBBSList(start, limit)
	c.JSON(200, resp)
}

// 管理 bbs
func HandleAdminBBSList(c *gin.Context) {
	user := getUser(c)
	start := c.DefaultQuery("start", "0")
	limit := c.DefaultQuery("limit", "10")

	resp := model.AdminBBSList(user, start, limit)
	c.JSON(200, resp)
}
