package common

import (
	"sync"
	"time"
)

// Cache 一个用户状态的缓存
type Cache struct {
	*uploadedIDSyncMap
	Users map[int]*MsgStatus
}

type uploadedIDSyncMap struct {
	lock       sync.Mutex
	uploadedID map[string]bool
}

// IsUploadedID 是否是oss已经存在的ID
func (u *uploadedIDSyncMap) IsUploadedID(id string) bool {
	u.lock.Lock()
	defer u.lock.Unlock()
	if _, ok := u.uploadedID[id]; ok {
		return true
	}
	return false
}

// AddUploadedID 添加已知id
func (u *uploadedIDSyncMap) AddUploadedID(id string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.uploadedID[id] = true
}

// AddUser 添加一个用户状态
func (c *Cache) AddUser(id int, m *MsgStatus) {
	c.Users[id] = m

}

// DeleteUser 删除用户
func (c *Cache) DeleteUser(id int) {
	delete(c.Users, id)
}

// AutoRemove 自动删除
func (c *Cache) AutoRemove() {
	for {
		time.Sleep(time.Second * 5)
		if len(c.Users) > 10000 {
			for _, i := range c.Users {
				if time.Since(i.time) > time.Hour*4 {
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
