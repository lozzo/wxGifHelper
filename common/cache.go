package common

import (
	"sync"
	"time"

	"github.com/golang/glog"
)

// Cache 一个用户状态的缓存
type Cache struct {
	UploadedIDSyncMap *uploadedIDSyncMap
	Users             map[int]*MsgStatus
}

type uploadedIDSyncMap struct {
	uploadedID map[string]bool
	lock       sync.Mutex
}

// Init 程序开始时，将所有sql内已有记录的数据放到已经加载的map内
func (c *Cache) Init(x []string) {
	if c.UploadedIDSyncMap == nil {
		c.UploadedIDSyncMap = &uploadedIDSyncMap{uploadedID: make(map[string]bool)}
	}
	for _, v := range x {
		c.UploadedIDSyncMap.uploadedID[v] = true
	}
	go c.autoRemove()
	// glog.V(5).Info(*c.UploadedIDSyncMap)
}

// IsUploadedID 是否是oss已经存在的ID
func (c *Cache) IsUploadedID(id string) bool {
	c.UploadedIDSyncMap.lock.Lock()
	defer c.UploadedIDSyncMap.lock.Unlock()
	if _, ok := c.UploadedIDSyncMap.uploadedID[id]; ok {
		return true
	}
	return false
}

// AddUploadedID 添加已知id
func (c *Cache) AddUploadedID(id string) {
	c.UploadedIDSyncMap.lock.Lock()
	defer c.UploadedIDSyncMap.lock.Unlock()
	c.UploadedIDSyncMap.uploadedID[id] = true
}

// AddUser 添加一个用户状态
func (c *Cache) AddUser(id int, m *MsgStatus) {
	c.Users[id] = m

}

// DeleteUser 删除用户
func (c *Cache) DeleteUser(id int) {
	delete(c.Users, id)
}

// autoRemove 自动删除
func (c *Cache) autoRemove() {
	for {
		time.Sleep(time.Second * 5)
		if len(c.Users) > 100000 {
			for _, i := range c.Users {
				if time.Since(i.time) > time.Hour*4 {
					glog.V(5).Info("清除", i.ID)
					c.DeleteUser(i.ID)
				}
			}
		}
	}
}

// GetUserMsgStatus 获取用户状态
func (c *Cache) GetUserMsgStatus(id int) *MsgStatus {
	if v, ok := c.Users[id]; ok {
		return v
	}
	m := &MsgStatus{time: time.Now(), BindWxCode: -1, ID: id}
	c.AddUser(id, m)
	return m
}

// SetUserMsgStatus 设置用户状态  -- 不用
func (c *Cache) SetUserMsgStatus(id int, m *MsgStatus) {
	c.AddUser(id, m)
}
