package bot

// bot包需要的通用工具
import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"tg_gif/common"
	"tg_gif/model"
	"tg_gif/tools"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang/glog"
)

// QUEUE 结束表情发送后处理数据队列的名称

// 判断是否绑定微信
func isBindWx(m *common.MsgStatus) (string, bool) {
	switch m.BindWxCode {
	case -1:
		t := model.TgUser{ID: m.ID}
		name, err := model.IsBindWx(&t)
		if err != nil {
			return "", false
		}
		if name == "" {
			m.BindWxCode = 0
			return "", false
		}
		m.BindWxCode = 1
		m.WxNickName = name
		glog.V(5).Info("+++", name)
		return name, true
	case 0:
		return "", false
	case 1:
		glog.V(5).Info(".11.", m.WxNickName)
		return m.WxNickName, true
	}
	return "", false
}

// 解绑微信
func unBindWx(m *common.MsgStatus) {
	t := model.TgUser{ID: m.ID}
	m.BindWxCode = 0
	model.UnBindWxFromTg(&t)
}

// StopSend 结束发送图片,从cache内删除，获取URL，转入队列
func StopSend(id int, m *common.MsgStatus) error {
	key := strconv.Itoa(id)
	data := m.JSON()
	if data == nil {
		str := fmt.Sprintf("%d个表情发送全军阵亡", len(m.File))
		remsg := tgbotapi.NewMessage(int64(m.ID), str)
		BotAPI.Send(remsg)
		return errors.New("")
	}
	err := tools.Enqueue(data, common.QUEUE)
	if err != nil {
		glog.Error("结束表情失败，数据入队失败：", err)
		return err
	}
	if tools.DelKey(key) != nil {
		glog.Error("删除ke失败：", err)
	}
	cache.DeleteUser(m.ID)
	return nil
}

func getFiles() {
	var files []*common.FileWithURL
	data, _ := tools.Dequeue(common.QUEUE)
	m := &common.MsgStatus{}
	if json.Unmarshal(data, m) != nil {
		// glog.Error("解析数据错误", err)
		return
	}

	for _, i := range m.File {
		if cache.IsUploadedID(i.ID) {
			i.URL = "uploaded"
		}
		file := common.FileWithURL{
			URL:  i.URL,
			Name: i.ID,
		}
		cache.AddUploadedID(i.ID)
		files = append(files, &file)
	}
	// todo 这儿应该是要将上传和写数据库放在同一个事物内进行,但是这儿没有,后面需要修复
	tools.DowAndUploadToOss(files, 10)
	model.AddFilesFromTg(m)
	str := fmt.Sprintf("%d个表情发送成功上传成功，请在微信小程序上查看", len(m.File))
	remsg := tgbotapi.NewMessage(int64(m.ID), str)
	BotAPI.Send(remsg)
}

// RUNDOW 开始下载
func RUNDOW() {
	for {
		getFiles()
		time.Sleep(time.Millisecond * 10)
	}
}
