package bot

import (
	"fmt"

	"github.com/golang/glog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// BotAPI 全局的BotAPI
var BotAPI *tgbotapi.BotAPI

// Conf TgBot配置参数
type Conf struct {
	Token      string `yaml:"token"`
	WebHookURL string `yaml:"webHookURL"`
}

// Init bot 初始化
func Init(c *Conf) {
	var err error
	BotAPI, err = tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		glog.Fatal(err)

	}
	res, err := BotAPI.SetWebhook(tgbotapi.NewWebhook(c.WebHookURL))
	if err != nil {
		glog.Fatal(err)
	}
	glog.V(5).Info(fmt.Sprintf("%s", string(res.Result)))
	me, _ := BotAPI.GetMe()
	glog.V(5).Info(fmt.Sprintf("%+v", me))
	go RUNDOW()
}

// SendText 给指定ID发送文字信息封装
func SendText(msg *tgbotapi.Message, text string) {
	remsg := tgbotapi.NewMessage(msg.Chat.ID, text)
	remsg.ReplyToMessageID = msg.MessageID
	BotAPI.Send(remsg)
}

// SendImage 根据fileID发送图片给用户   AgADBQADR6gxGyFykFRBggQGx_gOBPpr2zIABKXCbZdtJm3wI0wCAAEC
func SendImage(msg *tgbotapi.Message, fileID string) {
	x := tgbotapi.NewPhotoShare(msg.Chat.ID, fileID)
	BotAPI.Send(x)
}
