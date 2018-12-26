package bot

// 处理消息的地方
// 用一个状态机维护消息
import (
	"fmt"
	"tg_gif/common"
	"tg_gif/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	helpStr = `wxGifHelper 是一个帮助将Telegram表情包发送至微信的小工具，配合微信小程序使用，小程序扫描二维码，或者搜索小程序处{}即可添加小程序需要绑定微信号才能使用，绑定后即可使用以下命令，你的绑定号是%d

			可用命令包括以下：

			/send       开始发送表情包
			/stop        结束发送/start状态结束，此时后台才才会开始下载表情包并上传到oss，小程序内更新可能会有一定延迟，上传结束后会有Tg通知
			/bind_wx     绑定微信帐号，也可以在微信程序内绑定TG帐号
			/un_bind_wx  解除TG帐号和微信帐号的绑定，解除绑定后已上传的表情依然可以在小程序查看
			
			希望使用开心
			`
	bindWxStr      = "绑定微信,请在微信小程序设置页面绑定：%d ID,扫描二维码或者搜索小程序:wxGifHelper"
	bindWxErrStr   = "你已绑定微信，微信昵称：%s，请勿重复绑定，或者解除绑定后重新绑定"
	unbindWxStr    = "当前绑定微信昵称：%s（昵称更新可能有有延迟）"
	unbindWxErrStr = "你尚未绑定微信"
	startSendStr   = "请开始发送表情，结束发送后请输入结束命令：/stop 本次发送状态保持时长为4小时，发送新表情重新计算时长，超过时长自动结束"
	commandErrStr  = "当前状态不可接受 %s 命令。可接受命令为 %s"
	notBindWxStr   = "当前尚未绑定微信，无法使用，请在小程序设置页面绑定Telegram ID:%d"
	wxAppQR        = "" // 微信小程序app二维码
)
var cache = common.Cache{Users: make(map[int]*common.MsgStatus)}

// Msg webHook收到的消息类型，需要bot处理
type Msg struct {
	UpdateID int               `json:"update_id"`
	Message  *tgbotapi.Message `json:"message"`
}

func (m *Msg) Handler() {
	userMsgStatus := GetUserMsgStatus(m.Message.From.ID)
	if m.Message.Entities != nil {
		x := *m.Message.Entities
		if x[0].Type == "bot_command" {
			commndHandler(m, userMsgStatus)
			return
		}
		return
	} else if m.Message.Sticker != nil { // 当包含表情时的时候
		FileID := m.Message.Sticker.FileID
		x := tgbotapi.FileConfig{FileID: FileID}
		f, err := BotAPI.GetFile(x)
		if err != nil {
			remsg := tgbotapi.NewMessage(m.Message.Chat.ID, "本次表情发送失败，请重新发送")
			remsg.ReplyToMessageID = m.Message.MessageID
			BotAPI.Send(remsg)
		}
		fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", BotAPI.Token, f.FilePath)
		file := common.GifORMp4{ID: m.Message.Sticker.FileID, Type: "Sticker", URL: fileURL}
		code, lenth := userMsgStatus.AppendFile(file)
		switch code {
		case 0:
			SendText(m.Message, fmt.Sprintf("收到 %d 个表情", lenth))
		case 1:
			SendText(m.Message, "表情已存在，请勿重复发送")
		case 2:
			SendText(m.Message, "请发送 /send 命令后再开始发送表情")
		}
	}

}

func commndHandler(m *Msg, userMsgStatus *common.MsgStatus) {
	name, isBindOk := isBindWx(userMsgStatus)
	command := m.Message.Text
	data, ok := userMsgStatus.IsCmdAllowed(command)
	if !ok {
		SendText(m.Message, fmt.Sprintf(commandErrStr, command, data))
		return
	}
	switch command {
	case "/start":
		SendText(m.Message, fmt.Sprintf(helpStr, m.Message.From.ID))
		model.NewTgUser(m.Message.From.ID)
	case "/send":
		if !isBindOk {
			SendText(m.Message, fmt.Sprintf(notBindWxStr, m.Message.From.ID))
			return
		}
		userMsgStatus.Cmd = "/send"
		SetUserMsgStatus(m.Message.From.ID, userMsgStatus)
		SendText(m.Message, startSendStr)
		return

	case "/bind_wx":
		if isBindOk {
			SendText(m.Message, fmt.Sprintf(bindWxErrStr, name))
			return
		}
		SendText(m.Message, fmt.Sprintf(bindWxStr, m.Message.From.ID))
		SendImage(m.Message, wxAppQR)
		return
	case "/un_bind_wx":
		if isBindOk {
			SendText(m.Message, fmt.Sprintf(unbindWxStr, name))
			unBindWx(m.Message.From.ID)
			return
		}
		SendText(m.Message, fmt.Sprintf(notBindWxStr, m.Message.From.ID))
		return
	case "/stop":
		StopSend(m.Message.From.ID, userMsgStatus)
		if len(*userMsgStatus.File) == 0 {
			SendText(m.Message, "已发送0个表情。哟")
			return
		}
		stopSendStr := fmt.Sprintf("已发送 %d 个表情 ,请在用户 %s 用户内打开小程序查看", len(*userMsgStatus.File), name)
		SendText(m.Message, stopSendStr)
		return
	}
}

// GetUserMsgStatus 获取用户发图状态
func GetUserMsgStatus(id int) *common.MsgStatus {
	return cache.GetUserMsgStatus(id)
}

// SetUserMsgStatus 设置用户状态
func SetUserMsgStatus(id int, m *common.MsgStatus) {
	cache.AddUser(id, m)
}
