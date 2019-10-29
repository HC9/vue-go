package model

import "time"

type articleResponse struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

// 首页文章传输的内容
type indexArticleResponse struct {
	Code   int       `json:"code"`
	Employ []Article `json:"employ"`
	News   []Article `json:"news"`
}

// 文章列表
type listArticleResponse struct {
	Code  int       `json:"code"`
	Item  []Article `json:"item"`
	Total int       `json:"total"`
}
