package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"vgo/cache"
	"vgo/service"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// 用户表
type User struct {
	Id         int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT"` // 主键
	Username   string `gorm:"size:30;UNIQUE"`             // 姓名
	Password   string `gorm:"size:61"`                    // 密码
	Email      string `gorm:"size:30;UNIQUE"`             // 邮箱，唯一，一个邮箱只能注册一次
	Role       int    `gorm:"DEFAULT:1"`                  // 默认权限为 1， 0 为管理员
	CreateTime time.Time
	LoginTime  time.Time
	Status     string `gorm:"DEFAULT:'active';size:10"` // 账号状态  active 生效 inactive 未生效 suspend 封号状态
	Avatar     string `gorm:"DEFAULT:'avatar.jpg';size:12"`
}

// 查询用户名和邮箱是否已被注册
func (user *User) CheckNameAndEmail(register *service.UserRegisterService) *service.Response {
	resp := &service.Response{}

	var fieldValue = map[string]string{"email": register.Email, "username": register.UserName}
	var filed = []string{"email", "username"}
	for _, key := range filed {
		tempUser := User{}
		DB.Where(fmt.Sprintf("%s = ?", key), fieldValue[key]).First(&tempUser)
		if tempUser.Id != 0 {
			// 查到重复数据
			resp.Error = fmt.Sprintf("%s have registered", fieldValue[key])
			break
		}
	}

	if resp.Error != "" {
		resp.Code = 52001
		resp.Msg = "注册校验失败"
	} else {
		resp.Code = 20000
		resp.Msg = "注册校验成功"
	}
	return resp
}

// 匹配用户输入的信息是否正确
func (user *User) LoginCheck(login *service.UserLoginService) *service.Response {
	// login 的 Name 指代 学工号 id_number 和 邮箱，即用邮箱或者学工号都可以登录
	DB.Where("email = ?", login.Username).Or("username = ?", login.Username).Find(&user)
	if user.Id != 0 {
		if err := user.CheckPassword(login.Password); err == nil {
			// 更新登录时间
			DB.Model(&user).Update("login_time", time.Now())
			resp := &service.Response{}
			resp.Code = 20000
			resp.Msg = "登录成功"
			temp := make(map[string]string)

			if user.Role == 0 {
				temp["token"] = "admin"
			} else {
				temp["token"] = "user"
			}
			resp.Data = temp
			return resp
		}
	}
	return &service.Response{Code: 51003, Msg: "登录信息错误"}
}

// 创建用户
func (user *User) Create(register *service.UserRegisterService) {
	user.Username = register.UserName
	user.Email = register.Email
	user.setPassword(register.Password) // 加密密码
	user.Role = 1
	user.CreateTime = time.Now()
	user.LoginTime = time.Now() // 默认第一次登录时间为创建时间

	if err := DB.Create(&user).Error; err != nil {
		DB.Rollback()
	}
}

// 加密密码
func (user *User) setPassword(password string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = fmt.Sprintf("%s", bytes)
}

// 校对密码
func (user *User) CheckPassword(password string) error {
	// password 为登录密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}

// 获取一个用户信息
// 单个用户完整信息不存储在 session 中，故如果用户需要频繁使用该信息，则需要做缓存处理
// 缓存有效时间为1个小时，查询会先经过缓存，然后再向 Mysql 中查找
func GetOneUser(ID string) *User {
	userStr, _ := cache.RedisClient.Get(ID).Result()

	user := &User{}
	if userStr == "" {
		// 查询  mysql
		r, _ := strconv.Atoi(ID)
		user.Id = r
		DB.First(user)
		// 存储数据到 cache
		usJs, _ := json.Marshal(user)
		cache.RedisClient.Set(ID, usJs, 3600*time.Second)
	} else {
		_ = json.Unmarshal([]byte(userStr), user)
	}
	return user
}

// 更新缓存信息
func UpdateUserCache(user *User) {
	usJs, _ := json.Marshal(user)
	ID := strconv.Itoa(user.Id)
	cache.RedisClient.Set(ID, usJs, 3600*time.Second)
}

// 登录用户修改密码
func (user *User) AdminUpdatePassword(password string) *service.Response {
	user.setPassword(password)
	DB.Model(&user).Update("password", user.Password)
	return &service.Response{
		Code:  20000,
		Data:  nil,
		Msg:   "修改密码成功",
		Error: "",
	}
}

// 修改用户头像
func (user *User) HandleUpdateAvatar(avatarName *string) {
	if !strings.HasPrefix(user.Avatar, strconv.Itoa(user.Id)) {
		DB.Model(user).Update("avatar", *avatarName)
		user.Avatar = *avatarName
	}
}
