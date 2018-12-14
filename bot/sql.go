package bot

import (
	"strconv"
	"tg_gif/tools"
)

// 未完成的函数。。。。。
// 判断是否绑定微信
func isBindWx(id int) (string, bool) {
	return "lozzow", true
}

// 解绑微信
func unBindWx(id int) {

}

// StopSend 结束发送图片,删除redis内的用户状态，转送至队列处理,处理时，如有重复文件，直接引用，不用上传服务器
func StopSend(id int, m *MsgStatus) error {
	tools.DelKey(strconv.Itoa(id))

	return nil
}
