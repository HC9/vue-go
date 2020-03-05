package utils

import (
	"fmt"
	"math/rand"
	"time"
	"vgo/model"

	"golang.org/x/crypto/bcrypt"
)

// 获取随机字符串
// length 为获取字符串的长度
func GetRandString(length int) string {
	s := []byte("abcdefghijklnmopqrstuvwxyzABCDEFGHIJKLNMOPQRSTUVWSYZ0123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randStr []byte
	for i := 0; i < length; i++ {
		randStr = append(randStr, s[r.Intn(len(s))])
	}
	return fmt.Sprintf("%s", randStr)
}

func fillData(user *model.User) {
	user.Username = GetRandString(10)
	user.Email = fmt.Sprintf("%s@qq.com", GetRandString(10))
	pass, _ := bcrypt.GenerateFromPassword([]byte("admin888"), 12)
	user.Password = fmt.Sprintf("%s", pass)
	user.Role = 1
	user.CreatedAt = time.Now()
	user.Status = "active"

}

// 填充数据库
func DevFillDB() {
	user := model.User{}
	model.DB.Where(model.User{}).First(&user)
	if user.ID == 0 {
		for i := 0; i < 50; i++ {
			user = model.User{}
			fillData(&user)
			model.DB.Save(&user)
		}
	}
}
