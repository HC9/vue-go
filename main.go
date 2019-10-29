package main

import (
	// conf 负责初始化，顺序不能换
	"os"
	_ "vgo/cache"
	_ "vgo/conf"
	_ "vgo/model"
	"vgo/route"
)

func main() {
	// 关闭数据库

	engine := route.NewRouter()
	// 运行
	_ = engine.Run(os.Getenv("RUN_ADDRESS")) // listen and serve on 0.0.0.0:8080
}
