package main

//vlog 5 DEBUG
import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"tg_gif/bot"
	"tg_gif/model"
	"tg_gif/server"
	"tg_gif/tools"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	yaml "gopkg.in/yaml.v2"
)

//wx配置
type wxConf struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

// Conf 总配置文件
type Conf struct {
	Bot   *bot.Conf        `yaml:"bot"`
	Redis *tools.RedisConf `yaml:"redis"`
	WX    *wxConf          `yaml:"wx"`
	DB    *model.SQLConf   `yaml:"DB"`
	Oss   *tools.OssConf   `yaml:"oss"`
}

func main() {
	flag.Parse()
	glog.V(5).Info("Under DEBUG Mode!!")
	configFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		panic(fmt.Sprintf("读取文件错误: %s", err))
	}
	conf := new(Conf)
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		panic(fmt.Sprintf("解析文件错误: %s", err))
	}
	glog.V(5).Info(fmt.Sprintf("%+v", conf.Redis))

	url, err := url.Parse(conf.Bot.WebHookURL)
	if err != nil {
		panic(fmt.Sprintf("获取webHookURL错误: %s", err))
	}
	if len(url.Path) < 20 {
		glog.Warningln("webHookURL 长度小于20")
	}
	tools.OssInit(conf.Oss)
	tools.RedisPoolInit(conf.Redis)
	model.DBInit(conf.DB)
	bot.Init(conf.Bot)
	server.Init(url.Path)
	server.Run(":8889")
}
