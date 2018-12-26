package common

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/golang/glog"
)

// MsgStatus 当前消息状态
type MsgStatus struct {
	lock       sync.Mutex
	Cmd        string      `json:"-"`    // 当前命令
	File       *[]GifORMp4 `json:"File"` // 文件列表
	ID         int         `json:"-"`
	time       time.Time
	WxNickName string `json:"-"`
	BindWxCode int    `json:"-"` // -1:未知，需要查询，0：未绑定，1：已绑定
}

// GifORMp4 动图或者mp4
type GifORMp4 struct {
	ID   string //FileID
	Type string // gif or MP4
	URL  string
}

// FileWithURL 网络文件
type FileWithURL struct {
	URL  string
	Name string
}

const QUEUE = "data_queue"

// AppendFile 添加文件
// 状态码 0：ok，1：文件已存在，2：当前状态不可发送
func (m *MsgStatus) AppendFile(g GifORMp4) (int, int) { //状态码。长度
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.Cmd != "/send" {
		return 0, 2
	}
	for _, v := range *m.File {
		if g.ID == v.ID {
			return 0, 1
		}
	}
	*m.File = append(*m.File, g)
	return 0, len(*m.File)
}

// IsCmdAllowed 状态判断写的真垃圾啊，要重写要重写，要上状态机！！！
func (m *MsgStatus) IsCmdAllowed(cmd string) ([]string, bool) {
	a := false
	allCmd := []string{"/send", "/stop", "/bind_wx", "/un_bind_wx"} // 所有可选状态
	startStatus := []string{"/stop"}                                // 处于开始状态的时候
	noStatus := []string{"/send", "/bind_wx", "/un_bind_wx"}        // 没有状态的时候
	for _, x := range allCmd {
		a = (a || (x == cmd))
	}
	if !a {
		return noStatus, a
	}
	var cmdMap = map[string][]string{
		"/stop":       noStatus,
		"/send":       startStatus,
		"/bind_wx":    noStatus,
		"/un_bind_wx": noStatus,
	}
	if m.Cmd == "" {
		if cmd == "/stop_send" {
			return noStatus, false
		}
		return allCmd, true
	}
	allowedCmd := cmdMap[m.Cmd]
	b := false
	for _, x := range allowedCmd {
		b = b || (x == cmd)
	}
	return noStatus, b
}

// JSON 返回json字符串
func (m *MsgStatus) JSON() []byte {
	x, err := json.Marshal(m)
	if err != nil {
		glog.Error("序列化数据错误", err)
		return nil
	}
	return x

}
