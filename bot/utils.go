package bot

// bot包需要的通用工具
import (
	"strconv"
	"tg_gif/model"
	"tg_gif/tools"

	"github.com/golang/glog"
)

// QUEUE 结束表情发送后处理数据队列的名称
const QUEUE = "data_queue"

// 未完成的函数。。。。。
// 判断是否绑定微信
func isBindWx(id int) (string, bool) {
	t := model.TgUser{id}
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
	t := model.TgUser{id}
	model.UnBindWx(&t)
}

// StopSend 结束发送图片,删除redis内的用户状态，转送至队列处理,处理时，如有重复文件，直接引用，不用上传服务器
func StopSend(id int, m *MsgStatus) error {
	key := strconv.Itoa(id)
	data, err := tools.GetByteValue(key)
	if err != nil {
		glog.Error("结束表情失败，获取redis数据失败：", err)
		return err
	}
	err = tools.Enqueue(data, QUEUE)
	if err != nil {
		glog.Error("结束表情失败，数据入队失败：", err)
		return err
	}
	if tools.DelKey(key) != nil {
		glog.Error("删除ke失败：", err)
	}
	return nil
}
