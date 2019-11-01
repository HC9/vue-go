package model

import (
	"errors"
	"strconv"
	"time"
	"vgo/service"
)

type Article struct {
	ID         int       `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"` // 主键
	UID        int       `gorm:"type:int" json:"uid"`
	Username   string    `gorm:"size:15;" json:"username"`
	Subject    string    `gorm:"size:17;" json:"subject"`
	Title      string    `gorm:"size:70;" json:"title"`
	Content    string    `gorm:"type:longtext" json:"content"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	Status     string    `gorm:"DEFAULT:'active';size:10" json:"status"`
}

// 创建文章，文章的 UID 与 username 通过缓存中的用户信息来获取
// subject 分 news, employ
// 前三个需要管理员权限和指定管理者才能添加，分别为0号为2号权限
// 1 号权限仅能添加版本趣闻 interesting 的内容
func (article *Article) Create(user *User) *service.Response {
	article.UID = user.Id
	article.Username = user.Username
	article.CreateTime = time.Now()
	article.Status = "active"

	switch user.Role {
	case 0:
		if err := DB.Create(&article).Error; err != nil {
			DB.Rollback()
			return &service.Response{
				Code:  55000,
				Data:  nil,
				Msg:   "",
				Error: err.Error(),
			}
		}
	case 1:
		err := errors.New("无权限添加文章")
		return &service.Response{
			Code:  55001,
			Data:  nil,
			Msg:   "",
			Error: err.Error(),
		}
	}

	return &service.Response{
		Code:  20000,
		Data:  nil,
		Msg:   "添加文章成功",
		Error: "",
	}

}

// 删除文章，软删除
// 将状态设计为 inactive
func (article *Article) Delete(user *User) *service.Response {

	e := errors.New("")
	// 管理员可删除任何文章
	if user.Role == 0 {
		e = DB.Model(&article).Where("id = ?", article.ID).Update("status", "inactive").Error
	} else {
		e = DB.Model(&article).Where("uid = ? and id = ?", user.Id, article.ID).Update("status", "inactive").Error
	}

	if e == nil {
		return &service.Response{
			Code:  20000,
			Data:  nil,
			Msg:   "删除文章成功",
			Error: "",
		}
	} else {
		return &service.Response{
			Code:  55003,
			Data:  nil,
			Msg:   "删除文章失败",
			Error: e.Error(),
		}
	}
}

// 获取文章详情
// 根据ID号来获取
func GetArticleInfo(id string) *service.Response {
	artResp := &articleResponse{}
	err := DB.Table("articles").Where("id = ? and `status`='active'", id).Scan(artResp).Error
	if err != nil {
		return &service.Response{
			Code:  54002,
			Data:  nil,
			Msg:   "无此文章数据",
			Error: err.Error(),
		}
	} else {
		return &service.Response{
			Code:  20000,
			Data:  artResp,
			Msg:   "获取文章详情成功",
			Error: "",
		}
	}
}

// 获取首页的两个新闻内容
// 先获取数据库当前的类新闻总数
func GetArticleIndex() *indexArticleResponse {
	resp := &indexArticleResponse{}

	NewsCount, EmployCount := 0, 0
	DB.Limit(8).Table("articles").Order("create_time desc").Where("status='active' and subject='news'").Count(&NewsCount)
	DB.Limit(8).Table("articles").Order("create_time desc").Where("status='active' and subject='employ'").Count(&EmployCount)

	resp.News = make([]Article, NewsCount)
	resp.Employ = make([]Article, EmployCount)

	resp.News = BuildIndexArticleResponse(NewsCount, "news")
	resp.Employ = BuildIndexArticleResponse(EmployCount, "employ")

	resp.Code = 20000
	return resp
}

// 获取最新的新闻列表
func GetArticleNews() *listArticleResponse {
	resp := listArticleResponse{}
	DB.Limit(6).Table("articles").Where("status='active'").Order("create_time desc").Scan(&resp.Item)
	resp.Code = 20000
	return &resp
}

// 文章列表详情处理，即某一类文章的分页处理
func GetArticleList(subject, start, limit string) *listArticleResponse {
	startInt, _ := strconv.Atoi(start)
	limitInt, _ := strconv.Atoi(limit)
	resp := listArticleResponse{}

	DB.Model(&Article{}).Where("subject = ? and status='active'", subject).Count(&resp.Total)
	DB.Limit(limitInt).Offset(startInt).Model(&Article{}).Where("subject = ? and status='active'", subject).Find(&resp.Item)

	resp.Code = 20000
	return &resp
}

// Todo
// 文章管理
func AdminArticleList(user *User, subject, start, limit string) *listArticleResponse {
	// 权限 0 代表管理员
	startInt, _ := strconv.Atoi(start)
	limitInt, _ := strconv.Atoi(limit)
	resp := listArticleResponse{}
	queryDB := DB.Table("articles").Limit(limitInt).Offset(startInt).Order("create_time desc").Where("status='active'")

	// 管理员有权对所有文章进行修改
	// 普通用户仅可对本人的文章进行修改
	switch user.Role {
	case 0:
		// 管理员
		// 需不需要分类来做
		if subject != "" {
			DB.Model(&Article{}).Where("subject = ? and status='active'", subject).Count(&resp.Total)
			queryDB.Where("subject = ? and status='active'", subject).Scan(&resp.Item)
		} else {
			DB.Model(&Article{}).Where("status='active'").Count(&resp.Total)
			queryDB.Where("status='active'").Scan(&resp.Item)
		}
	case 1:
		// 普通用户
		if subject != "" {
			DB.Model(&Article{}).Where("subject = ? and uid= ? and status='active'", subject, user.Id).Count(&resp.Total)
			queryDB.Where("subject = ? and uid= ? and status='active'", subject, user.Id).Scan(&resp.Item)
		} else {
			DB.Model(&Article{}).Where("uid = ? and status='active'", user.Id).Count(&resp.Total)
			queryDB.Where("uid = ? and status='active'", user.Id).Scan(&resp.Item)
		}
	}

	resp.Code = 20000
	return &resp
}

// 更新文章内容
// 可以更新文章的标题与内容，添加更新时间
func UpdateArticle(user *User, id int, title, content string) *service.Response {

	article := Article{ID: id, UID: user.Id}
	DB.First(&article)
	// 更新内容
	article.Content = content
	article.Title = title
	article.UpdateTime = time.Now()

	e := DB.Model(&article).Updates(Article{Content: article.Content, Title: article.Title, UpdateTime: article.UpdateTime}).Error
	if e == nil {
		return &service.Response{
			Code:  20000,
			Data:  nil,
			Msg:   "更新文章成功",
			Error: "",
		}
	} else {
		return &service.Response{
			Code:  55004,
			Data:  nil,
			Msg:   "删除文章失败",
			Error: "无权限删除或无指定文章",
		}
	}
}
