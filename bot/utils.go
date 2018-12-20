package bot

// bot包需要的通用工具
import (
	"encoding/json"
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

// 未完成的函数。。。。。
// 判断是否绑定微信
func isBindWx(id int) (string, bool) {
	t := model.TgUser{ID: id}
	name, err := model.IsBindWx(&t)
	if err != nil {
		return "", false
	}
	if name == "" {
		return "", false
	}
	return name, true
}

// 解绑微信
func unBindWx(id int) {
	t := model.TgUser{ID: id}
	model.UnBindWx(&t)
}

// StopSend 结束发送图片,删除redis内的用户状态，转送至队列处理,处理时，如有重复文件，直接引用，不用上传服务器
func StopSend(id int, m *common.MsgStatus) error {
	var err error
	key := strconv.Itoa(id)
	if m.Count == 0 {
		if tools.DelKey(key) != nil {
			glog.Error("删除ke失败：", err)
			return err
		}
		return nil
	}
	data, err := tools.GetByteValue(key)
	if err != nil {
		glog.Error("结束表情失败，获取redis数据失败：", err)
		return err
	}
	err = tools.Enqueue(data, common.QUEUE)
	if err != nil {
		glog.Error("结束表情失败，数据入队失败：", err)
		return err
	}
	if tools.DelKey(key) != nil {
		glog.Error("删除ke失败：", err)
	}
	// GetFiles()
	// model.AddFilesFromTg(m)
	return nil
}

func GetFiles() {
	var files []*common.FileWithURL
	data, _ := tools.Dequeue(common.QUEUE)
	m := &common.MsgStatus{}
	if json.Unmarshal(data, m) != nil {
		// glog.Error("解析数据错误", err)
		return
	}

	for _, i := range *m.File {
		file := common.FileWithURL{
			URL:  i.URL,
			Name: i.ID,
		}
		files = append(files, &file)
	}
	tools.DowAndUploadToOss(files, 10)
	model.AddFilesFromTg(m)
	str := fmt.Sprintf("%d个表情发送成功上传成功，请在微信小程序上查看", m.Count)
	if m.IsGroup {
		str = fmt.Sprintf(",表情包组：%s,共%d个表情发送成功上传成功，请在微信小程序上查看", m.GroupName, m.Count)
	}
	remsg := tgbotapi.NewMessage(int64(m.ID), str)
	BotAPI.Send(remsg)
}

// RUNDOW 开始下载
func RUNDOW() {
	for {
		GetFiles()
		time.Sleep(time.Microsecond * 10)
	}
}
