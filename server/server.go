package server

import (
	"encoding/json"
	"io/ioutil"
	"tg_gif/bot"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var r *gin.Engine

// Run http服务启动
func Run(ipPort string) {
	r.Run(ipPort)
}

// Init webHook服务器
func Init(WebHookURL string) {
	r = gin.New()
	r.Use(gin.Recovery())
	r.POST(WebHookURL, baseHandler)
}
func baseHandler(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var msg bot.Msg
	glog.V(5).Info(string(x))
	// err := c.ShouldBindJSON(&msg)  不知为啥绑定不上
	err := json.Unmarshal(x, &msg)
	if err != nil {
		glog.V(5).Info("bind err:", err)
	}
	msg.Handler()
}
