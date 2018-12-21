package model

import (
	"database/sql"
	"fmt"

	"github.com/golang/glog"
)

// 用户相关功能

// AddUser 添加用户，新用户，来源可能是wx，也可能是tg
func AddUser(u User) error {
	if u.TgUser != nil && u.WxUser == nil {
		_, err := newTgUser(u.TgUser)
		return err
	}
	if u.TgUser == nil && u.WxUser != nil {
		return newUserFromWx(u.WxUser)
	}
	return nil
}

// 新用户来自wx
func newUserFromWx(w *WxUser) error {
	userID, err := newWxUser(w)
	if err != nil {
		return err
	}
	stmt2, err := db.Prepare(`INSERT INTO users (wxUserID) VALUES (?)`)
	if err != nil {
		glog.Error("数据库错误：", err)
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(userID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return err
	}

	return nil
}

//添加新的Telegram用户
func newTgUser(t *TgUser) (int64, error) {
	stmt, err := db.Prepare(`INSERT  INTO tgUsers (id) SELECT (?) FROM DUAL WHERE NOT EXISTS (SELECT id FROM tgUsers WHERE id= ? )`)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID, t.ID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	return int64(t.ID), nil
}

// 添加新的Wx用户
func newWxUser(w *WxUser) (int64, error) {
	// stmt, err := db.Prepare(`INSERT INTO wxUsers (openID,nickName) VALUES (?,?)`)
	stmt, err := db.Prepare(`INSERT INTO wxUsers (openID,nickName) SELECT (?,?) FROM DUAL WHERE NOT EXISTS (SELECT openID FROM wxUsers WHERE openID= ? ) `)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(w.openID, w.NickName, w.openID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	return userID, nil
}

// BindTg 已有用户绑定新Tg
func BindTg(w *WxUser, t *TgUser) error {
	_, err := newTgUser(t)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(`UPDATE users SET tgUserID = ? WHERE wxUserID =（SELECT id FROM wxUsers WHERE openID = ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID, w.openID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return err
	}
	return nil
}

// UnBindWxFromTg 解绑Wx
func UnBindWxFromTg(t *TgUser) error {
	stmt, err := db.Prepare(`UPDATE users SET tgUserID = NULL WHERE tgUserID = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return err
	}
	return nil
}

// UnBindTgFromWX 解绑Tg
func UnBindTgFromWX(w *WxUser) error {
	stmt, err := db.Prepare(`UPDATE users SET tgUserID = NULL WHERE wxUserId = (SELECT id FROM wxUsers WHERE openID = ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.openID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return err
	}
	return nil
}

// IsBindWx 在tg内查看是否绑定wx,如果绑定则返回用户名，wx用户名不能是空字符串
func IsBindWx(t *TgUser) (string, error) {
	var u string
	err := db.QueryRow("SELECT nickName FROM wxUsers WHERE id = (SELECT wxUserID FROM users WHERE tgUserID = ? LIMIT 1) LIMIT 1 ", t.ID).Scan(&u)
	if err == sql.ErrNoRows {
		glog.V(2).Info(fmt.Sprintf("TG ID:%d 尚未绑定wx帐号", t.ID))
		return "", nil //这个时候是没有绑定
	}
	if err != nil {
		glog.Error("数据库错误：", err)
		return "", err //出现错误
	}
	return u, nil
}

// IsBindTg 在wx小程序内查产是否绑定tg帐号
func IsBindTg(w *WxUser) (int, error) {
	var t TgUser
	err := db.QueryRow("SELECT tgUserID FROM users WHERE wxUserID = (SELECT id FROM wxUsers WHERE openID = ?)", w.openID).Scan(&t.ID)
	if err == sql.ErrNoRows {
		glog.V(2).Info(fmt.Sprintf("微信openID:%s 用户尚未绑定TG帐号", w.openID))
		return 0, nil //这个时候是没有绑定
	}
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err //出现错误
	}
	return t.ID, nil
}
