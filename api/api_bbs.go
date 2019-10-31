package api

import (
	"vgo/model"
	"vgo/service"

	"github.com/gin-gonic/gin"
)

// 创建微博
func HandleCreateBBS(c *gin.Context) {
	user, _ := c.Get("user")
	if u, ok := user.(*model.User); ok {
		bbs := model.BBS{}
		bbs.Avatar = "/api/v1/user/avatar/" + u.Avatar
		_ = c.ShouldBindJSON(&bbs)
		bbs.Username = u.Username
		resp := bbs.Create()
		c.JSON(200, resp)
	}
}

// 添加评论
func HandleAddComment(c *gin.Context) {
	user, _ := c.Get("user")
	if u, ok := user.(*model.User); ok {
		postComment := &model.PostComment{}
		postComment.Username = u.Username
		_ = c.BindJSON(postComment)
		postComment.AddComment(u.Avatar)
	}
	c.JSON(200, gin.H{
		"code": 20000,
		"msg":  "评论添加成功",
	})
}

// 删除微博
func HandleDeleteBBS(c *gin.Context) {
	user, _ := c.Get("user")
	if u, ok := user.(*model.User); ok {
		w := &model.BBS{}
		_ = c.ShouldBindJSON(w)
		// 0 号权限是管理员，能够不用验证执行删帖操作
		resp := &service.Response{}
		if u.Role == 0 {
			resp = w.Delete()
		} else {
			// 对用户名进行验证
			info := w.GetBBSInfo(w.ID)
			if temp, ok := info.Data.(*model.BBS); ok {
				if temp.Username == u.Username {
					resp = w.Delete()
				} else {
					resp.Code = 55003
					resp.Msg = "你没有权限执行此操作"
				}
			}
		}

		c.JSON(200, resp)
	}
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
	user, _ := c.Get("user")
	if u, ok := user.(*model.User); ok {
		start := c.DefaultQuery("start", "0")
		limit := c.DefaultQuery("limit", "10")

		resp := model.AdminBBSList(u, start, limit)
		c.JSON(200, resp)
	}
}
