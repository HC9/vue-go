package conf

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/joho/godotenv"
)

// 将  zh-cn.yaml 翻译成结构体，供 json 验证使用
type TransCn struct {
	// 选了  Tag 标准的内容
	ActualTag struct {
		Required string `yaml:"required"`
		Min      string `yaml:"min"`
		Max      string `yaml:"max"`
		Email    string `yaml:"email"`
	} `yaml:"ActualTag"`

	Field struct {
		UserName        string `yaml:"UserName"`
		PassWord        string `yaml:"PassWord"`
		ConfirmPassword string `yaml:"ConfirmPassword"`
		ID              string `yaml:"ID"`
		Email           string `yaml:"Email"`
		Sex             string `yaml:"Sex,flow"`
	} `yaml:"Field"`
}

//var Dictionary *TransCn

var Dictionary map[string]map[string]string

// 提取配置文件进入系统变量
// 初始化包的，由编译器自动运行
func init() {
	// godotenv.Load 挂载项目根目录下的 .env 文件
	if err := godotenv.Load(); err != nil {
		// 没有配置文件，退出启动
		log.Fatal(err)
	}

	// 挂载 yaml 文件
	data, err := ioutil.ReadFile(os.Getenv("YAML_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	dict := make(map[string]map[string]string)
	//dict := TransCn{}
	if err := yaml.Unmarshal(data, dict); err != nil {
		log.Fatal(err)
	}
	// 全局 Dictionary
	Dictionary = dict

}
