package bot

// 用一个状态机维护消息
import (
	"encoding/json"
	"fmt"
	"strconv"

	"tg_gif/tools"

	"github.com/garyburd/redigo/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
)

var (
	helpStr        = "wxGifHelper 是一个帮助将Telegram表情包发送至微信的小工具，配合微信小程序使用，小程序扫描二维码，或者搜索小程序处{}即可添加小程序"
	bindWxStr      = "绑定微信,请在微信小程序设置页面绑定：%d ID,扫描二维码或者搜索小程序:wxGifHelper"
	bindWxErrStr   = "你已绑定微信，微信昵称：%s，请勿重复绑定，或者解除绑定后重新绑定"
	unbindWxStr    = "当前绑定微信昵称：%s，解除微信绑定后，你亦可以用自己的电报ID: %d 来查找已上传的表情"
	unbindWxErrStr = "你尚未绑定微信"
	startSendStr   = "请开始发送表情，结束发送后请输入结束命令：/stop_send 本次发送状态保持时长为1天"
	startGroupStr  = "请为当前表情包组命名"
	errorStr       = "当前状态为 %s"
	commandErrStr  = "当前状态不可接受 %s 命令。可接受命令为 %s"

	wxAppQR = "" // 微信小程序app二维码
)

// Msg webHook收到的消息类型，需要bot处理
type Msg struct {
	UpdateID int               `json:"update_id"`
	Message  *tgbotapi.Message `json:"message"`
}

// Handler 数据处理
func (m *Msg) Handler() {
	userMsgStatus := GetUserMsgStatus(m.Message.From.ID)
	if m.Message.Entities != nil {
		x := *m.Message.Entities
		if x[0].Type == "bot_command" {
			commndHandler(m, userMsgStatus)
			return
		} else {
			return
		}
	} else if m.Message.Sticker != nil { // 当包含表情时的时候
		if userMsgStatus.Cmd == "/start_send" || (userMsgStatus.Cmd == "/start_group" && userMsgStatus.Status == 1) {
			file := GifORMp4{m.Message.Sticker.FileID, "Sticker"}
			if userMsgStatus.File == nil {
				userMsgStatus.File = &[]GifORMp4{}
			}
			for _, exFile := range *userMsgStatus.File {
				if file.ID == exFile.ID {
					SendText(m.Message, "此表情本次已发送，请勿重复发送")
					return
				}
			}
			*userMsgStatus.File = append(*userMsgStatus.File, file)
			userMsgStatus.Count++
			SetUserMsgStatus(m.Message.From.ID, userMsgStatus)
			SendText(m.Message, fmt.Sprintf("收到 %d 个表情", userMsgStatus.Count))
			return
		}
		SendText(m.Message, "")

	} else if m.Message.Sticker == nil && m.Message.Entities == nil { // 当只含有文字信息时
		if userMsgStatus.Cmd == "/start_group" && userMsgStatus.Status == 0 { //此时等待用户输入表情包组名字，不能是空格
			if m.Message.Text != "" {
				userMsgStatus.Status = 1
				userMsgStatus.GroupName = m.Message.Text
				SetUserMsgStatus(m.Message.From.ID, userMsgStatus)
				SendText(m.Message, fmt.Sprintf("当前表情包组名为： %s,请开始发送表情包，/stop_send 结束发送", m.Message.Text))
				return
			} else {
				SendText(m.Message, "组名为空，请重新输入")
				return
			}
			
		}

	}
}

// 只是针对命令行的处理函数
func commndHandler(m *Msg, userMsgStatus *MsgStatus) {
	command := m.Message.Text
	data, ok := userMsgStatus.isCmdAllowed(command)
	if !ok {
		SendText(m.Message, fmt.Sprintf(commandErrStr, command, data))
		return
	}
	switch command {
	case "/bind_wx":
		name, ok := isBindWx(m.Message.From.ID)
		if ok {
			SendText(m.Message, fmt.Sprintf(bindWxErrStr, name))
			return
		}
		SendText(m.Message, fmt.Sprintf(bindWxStr, m.Message.From.ID))
		SendImage(m.Message, wxAppQR)
		return
	case "/un_bind_wx":
		name, ok := isBindWx(m.Message.From.ID)
		if ok {
			SendText(m.Message, fmt.Sprintf(unbindWxStr, name, m.Message.From.ID))
			unBindWx(m.Message.From.ID)
			return
		}
		SendText(m.Message, unbindWxErrStr)
		return
	case "/start_send":
		userMsgStatus.Cmd = "/start_send"
		userMsgStatus.Status = 1
		userMsgStatus.IsGroup = false
		SetUserMsgStatus(m.Message.From.ID, userMsgStatus)
		SendText(m.Message, startSendStr)
		return
	case "/start_group":
		userMsgStatus.Cmd = "/start_group"
		userMsgStatus.Status = 0
		userMsgStatus.IsGroup = true
		err := SetUserMsgStatus(m.Message.From.ID, userMsgStatus)
		x := GetUserMsgStatus(m.Message.From.ID)
		glog.V(5).Info(x, err)
		SendText(m.Message, startGroupStr)
		return
	case "/stop_send":
		StopSend(m.Message.From.ID, userMsgStatus)
		if userMsgStatus.Count == 0 {
			SendText(m.Message, "已发送0个表情。哟")
			return
		}
		stopSendStr := ""
		if userMsgStatus.IsGroup {
			stopSendStr = fmt.Sprintf("结束当前发送表情 %d个，组名为 %s ", userMsgStatus.Count, userMsgStatus.GroupName)
		} else {
			stopSendStr = fmt.Sprintf("结束当前发送表情 %d个", userMsgStatus.Count)
		}
		name, ok := isBindWx(m.Message.From.ID)
		if ok {
			stopSendStr = fmt.Sprintf("%s ,请在用户 %s用户内打开小程序查看", stopSendStr, name)
		} else {
			stopSendStr = fmt.Sprintf("%s ,请使用ID: %d 搜索查看", stopSendStr, m.Message.From.ID)
		}
		SendText(m.Message, stopSendStr)
		return
	}
}

// ALL 不同的用户有各自的状态。存在redis，避免自己写过期设置，
// 开始发送表情之后1小时都可以继续发送，超过一小时，删除发送状态，每次发图重新计算时间
// redis json 数据格式类型如下，每次更新重写设置时间
// var ALL = make(map[string]MsgStatus)

// MsgStatus 当前消息状态
type MsgStatus struct {
	Cmd       string      // 当前命令
	Count     int         // 图片数量
	File      *[]GifORMp4 // 文件列表
	IsGroup   bool        // 是否为一组文件
	Status    int         // 0 ：未开始存图。1：正在存图，2：结束存图 主要用在group时命名等待
	GroupName string
}

// GifORMp4 动图或者mp4
type GifORMp4 struct {
	ID   string //FileID
	Type string // gif or MP4
}

func (m *MsgStatus) appendFile(g GifORMp4) {
	if m.Status != 2 {
		*m.File = append(*m.File, g)
	}
}

// 状态判断写的真垃圾啊，要重写要重写，要上状态机！！！
func (m *MsgStatus) isCmdAllowed(cmd string) ([]string, bool) {
	a := false
	allCmd := []string{"/start_send", "/start_group", "/bind_wx", "/un_bind_wx", "/stop_send"}
	for _, x := range allCmd {
		a = a || (x == cmd)
	}
	if !a {
		return allCmd, a
	}
	cmds1 := []string{"/stop_send"}
	cmds2 := []string{"/start_send", "/start_group", "/bind_wx", "/un_bind_wx"}
	var x = map[string][]string{
		"/stop_send":   cmds2,
		"/start_send":  cmds1,
		"/start_group": cmds1,
		"/bind_wx":     cmds1,
		"/un_bind_wx":  cmds1,
	}
	if m.Cmd == "" {
		if cmd == "/stop_send" {
			return cmds2, false
		}
		return cmds2, true
	}
	allowedCmd := x[m.Cmd]
	b := false
	for _, x := range allowedCmd {
		b = b || (x == cmd)
	}
	return allowedCmd, b
}

// GetUserMsgStatus 获取用户发图状态
func GetUserMsgStatus(id int) *MsgStatus {
	m := &MsgStatus{}
	idString := strconv.Itoa(id)
	r, err := tools.GetByteValue(idString)
	if err == redis.ErrNil {
		glog.V(2).Info("该用户目前没维护状态，设置一个新的状态")
		return m
	} else if err != nil {
		glog.Error("redis错误：", err)
		return nil
	}

	if json.Unmarshal(r, m) != nil {
		fmt.Println(r, *m)
		glog.Error("解析数据错误", err)
		return nil
	}
	fmt.Printf("%+v", m)
	return m
}

// SetUserMsgStatus 设置用户状态
func SetUserMsgStatus(id int, m *MsgStatus) error {
	idString := strconv.Itoa(id)
	x, err := json.Marshal(m)
	if err != nil {
		glog.Error("序列化数据错误", err)
		return err
	}
	err = tools.SetValue(idString, string(x), 86400)
	if err != nil {
		glog.Error("redis错误：", err)
	}
	return err
}
