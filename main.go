package main

import (
	"os"
	_ "vgo/cache"
	_ "vgo/model"
	"vgo/route"
	_ "vgo/service"

	"github.com/joho/godotenv"
)

func main() {
	// 关闭数据库
	_ = godotenv.Load("D:\\GoWorkDir\\go-vgo\\.env")
	engine := route.NewRouter()
	// 运行
	_ = engine.Run(os.Getenv("RUN_ADDRESS")) // listen and serve on 0.0.0.0:8080
}
