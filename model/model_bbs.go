package model

import (
	"strconv"
	"time"
	"vgo/service"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

// 使用 mongodb 数据库记录每一条博文
/*
需要的数据
username 需要根据 cookies 提供的 user_id 向缓存中提取 username
id 自动生成
username 当前月户的名字，唯一
header 标题
content 内容
create_time
comment 评论，拥有以下字段
	username
	comment
	create_time
*/

// 评论结构体
type Comment struct {
	Username  string `json:"username" bson:"username"`
	Content   string `json:"content" bson:"content"`
	CreatTime int64  `json:"create_time" bson:"create_time"`
	Avatar    string `json:"avatar" bson:"avatar"`
}

// 微博主结构体
type BBS struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Username    string        `json:"username" bson:"username"`
	Title       string        `json:"title" bson:"title"`
	Content     string        `json:"content" bson:"content"`
	CreatTime   int64         `json:"create_time" bson:"create_time"`
	Avatar      string        `json:"avatar" bson:"avatar"`
	CommentList []Comment     `json:"comment" bson:"comment"`
}

// 创建 webibo
func (w *BBS) Create() *service.Response {
	w.ID = bson.NewObjectId()
	w.CommentList = make([]Comment, 1)
	w.CreatTime = time.Now().Unix()
	err := BBSCollection.Insert(w)

	if err != nil {
		return &service.Response{
			Code:  55001,
			Data:  nil,
			Msg:   "创建微博失败",
			Error: err.Error(),
		}
	} else {
		return &service.Response{
			Code:  20000,
			Data:  w,
			Msg:   "创建微博成功",
			Error: "",
		}
	}
}

// 删除微博
// 根据post的ID进行删除操作，还需要匹配用户的
func (w *BBS) Delete() *service.Response {
	_ = BBSCollection.Remove(bson.M{"_id": w.ID})
	return &service.Response{
		Code:  20000,
		Data:  nil,
		Msg:   "删除成功",
		Error: "",
	}
}

// 微博详情
// 先将 id 转换成 objectID(convertStrToID)，再使用FindOne查询
func (w *BBS) GetBBSInfo(id interface{}) *service.Response {
	if sID, ok := id.(string); ok {
		w.ID = bson.ObjectIdHex(sID)
	} else if pID, ok := id.(bson.ObjectId); ok {
		w.ID = pID
	}
	err := BBSCollection.Find(bson.M{"_id": w.ID}).One(w)
	if err != nil {
		return &service.Response{
			Code:  54002,
			Data:  nil,
			Msg:   "无此数据",
			Error: err.Error(),
		}
	} else {
		return &service.Response{
			Code:  20000,
			Data:  w,
			Msg:   "获取数据详情成功",
			Error: "",
		}
	}
}

// 增加评论
// ajax 提交的结构体
type PostComment struct {
	ID       bson.ObjectId `json:"id" bson:"id"`
	Username string        `json:"username" bson:"username"`
	Content  string        `json:"content" bson:"content"`
}

func (post *PostComment) AddComment(avatar string) {
	//comment := make(map[int64]map[string]string)
	//keyCreateTime := time.Now().Unix()
	//
	//thisComment := make(map[string]string)
	//thisComment["username"] = post.Username
	//thisComment["content"] = post.Content
	//
	//comment[keyCreateTime] = thisComment
	comment := Comment{}
	comment.Content = post.Content
	comment.Username = post.Username
	comment.CreatTime = time.Now().Unix()
	comment.Avatar = "/api/v1/user/avatar/" + avatar

	// addToSet 往数据组中进行追加内容
	// 参考文档 : https://blog.csdn.net/qq_16313365/article/details/58599253
	update := bson.D{
		{"$addToSet", bson.D{
			{"comment", comment},
		}},
	}
	_ = BBSCollection.Update(bson.M{"_id": post.ID}, update)
}

// ====================================================
type bbsListResponse struct {
	Total int                `json:"total"`
	Item  []bbsListSerialize `json:"item"`
	Code  int                `json:"code"`
}

type bbsListSerialize struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	Title     string        `json:"title" bson:"title"`
	CreatTime int64         `json:"create_time" bson:"create_time"`
}

// 获取 bbs 列表
func GetBBSList(start, limit string) *bbsListResponse {
	startInt, _ := strconv.Atoi(start)
	limitInt, _ := strconv.Atoi(limit)
	bbsListResp := bbsListResponse{}

	query := BBSCollection.Find(bson.M{}).Sort("-create_time")
	bbsListResp.Total, _ = query.Count()
	query = query.Skip(startInt).Limit(limitInt)

	respCount, _ := query.Count()
	bbsListResp.Item = make([]bbsListSerialize, respCount)
	_ = query.All(&bbsListResp.Item)
	bbsListResp.Code = 20000
	return &bbsListResp
}

// BBS 管理
func AdminBBSList(user *User, start, limit string) *bbsListResponse {
	// 权限 0 代表管理员
	startInt, _ := strconv.Atoi(start)
	limitInt, _ := strconv.Atoi(limit)
	resp := &bbsListResponse{}
	query := &mgo.Query{}
	// 管理员有权对所有文章进行修改
	// 普通用户仅可对本人的文章进行修改
	switch user.Role {
	case 0:
		query = BBSCollection.Find(bson.M{}).Sort("-create_time")
	case 1:
		query = BBSCollection.Find(bson.M{"username": user.Username}).Sort("-create_time")
	}

	resp.Total, _ = query.Count()
	query = query.Skip(startInt).Limit(limitInt)
	respCount, _ := query.Count()

	resp.Item = make([]bbsListSerialize, respCount)
	_ = query.All(&resp.Item)
	resp.Code = 20000
	return resp
}
