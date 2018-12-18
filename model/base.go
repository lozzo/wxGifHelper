package model

import (
	"database/sql"
	"fmt"
	"time"
)

//DB 全局的 数据库链接
var db *sql.DB

// SQLConf sql配置项
type SQLConf struct {
	DBUrl           string `yaml:"DBUrl"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

// Gifs gifs表映射
type Gifs struct {
	ID    int       `json:"id"`
	Group *GifGroup `json:"group"`
	File  string    `json:"fileID"`
	User  *User     `json:"user"`
}

// GifGroup gifGroups表
type GifGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// User users表
type User struct {
	ID     int     `json:"id"`
	TgUser *TgUser `json:"tg"`
	WxUser *WxUser `json:"wx"`
}

// TgUser tgUsers表
type TgUser struct {
	ID int `json:"id"`
}

// WxUser wxUsers表
type WxUser struct {
	ID       int    `json:"id"`
	NickName string `json:"nickName"`
	openID   string
}

// DBInit 数据库链接初始化
func DBInit(c *SQLConf) {
	var err error
	db, err = sql.Open("mysql", c.DBUrl)
	if err != nil {
		fmt.Println(err)
		panic("error while init sql")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic("error while init sql")
	}
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(c.ConnMaxLifetime))
}
