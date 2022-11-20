package conf

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var c *App

func init() {
	c = new(App)
	c = loadConfig()
}

// Conf 服务配置
type Conf struct {
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
	RecordId  uint64 `yaml:"record_id"`
	Domain    string `yaml:"domain"`
	SubDomain string `yaml:"sub_domain"`
	Time      int64  `yaml:"time"`
}

type App struct {
	Conf *Conf `yaml:"app"`
}

func Get() *Conf {
	return c.Conf
}

// 加载 app.yaml 配置
func loadConfig() *App {
	path := flag.String("c", "./app.yaml", "输入 -c xxx.yaml 自定义配置文件")
	flag.Parse()
	file, e := os.ReadFile(*path)
	if e != nil {
		panic(e)
	}

	var app App
	e = yaml.Unmarshal(file, &app)
	if e != nil {
		panic(e)
	}
	fmt.Println("🔨 Config -> " + *path + "\n")
	return &app
}
