package service

// 序列化响应结构
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

// 用户模型的响应结构体
type UserResponse struct {
	Username   string `json:"用户名"`
	Email      string `json:"邮箱"`
	CreateTime string `json:"创建时间"`
	LoginTime  string `json:"登录时间"`
	Status     string `json:"状态"`
}
