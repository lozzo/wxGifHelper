package common

// MsgStatus 当前消息状态
type MsgStatus struct {
	Cmd       string      // 当前命令
	Count     int         // 图片数量
	File      *[]GifORMp4 // 文件列表
	IsGroup   bool        // 是否为一组文件
	Status    int         // 0 ：未开始存图。1：正在存图，2：结束存图 主要用在group时命名等待
	GroupName string
	ID        int
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

// ALL 不同的用户有各自的状态。存在redis，避免自己写过期设置，
// 开始发送表情之后1小时都可以继续发送，超过一小时，删除发送状态，每次发图重新计算时间
// redis json 数据格式类型如下，每次更新重写设置时间
// var ALL = make(map[string]MsgStatus)

func (m *MsgStatus) AppendFile(g GifORMp4) {
	if m.Status != 2 {
		*m.File = append(*m.File, g)
	}
}

// 状态判断写的真垃圾啊，要重写要重写，要上状态机！！！
func (m *MsgStatus) IsCmdAllowed(cmd string) ([]string, bool) {
	a := false
	allCmd := []string{"/start_send", "/start_group", "/bind_wx", "/un_bind_wx", "/stop_send"}
	cmds1 := []string{"/stop_send"}
	cmds2 := []string{"/start_send", "/start_group", "/bind_wx", "/un_bind_wx"}
	if cmd == "/start" {
		return cmds2, true
	}
	for _, x := range allCmd {
		a = a || (x == cmd)
	}
	if !a {
		return cmds2, a
	}
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
	return cmds2, b
}
