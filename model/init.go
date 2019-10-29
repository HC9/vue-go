package model

import (
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB                       // mysql 数据库连接
var BBSCollection *mgo.Collection     // mongodb weibo 数据库连接
var ArticleCollection *mgo.Collection // mongodb article 数据库连接

// gorm 使用下划线来分隔驼峰命名法，且其自动名字转成小写
// CreatTime creat_time
// 初始化 model 并建立数据库表
func init() {
	connectMysql()
	connectMongo()
}

func connectMysql() {
	db, err := gorm.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	// 是否打印数据库信息
	if gin.Mode() == "release" {
		// 生产模式不打印
		db.LogMode(false)
	} else {
		db.LogMode(true)
	}
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(20)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)

	// 创建各表
	DB = db
	// 创表
	DB.Set("gorm:table_options", "charset=utf8mb4")
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Article{})
}

func connectMongo() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	ArticleCollection = session.DB("VGO").C("articles")
	BBSCollection = session.DB("VGO").C("bbs")
}
