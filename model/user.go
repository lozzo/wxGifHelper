package model

import (
	"database/sql"
	"fmt"

	"github.com/golang/glog"
)

// NewUserFromWx 新用户来自wx
func NewUserFromWx(w *WxUser) (int, error) {
	userID, err := newWxUser(w)
	if err != nil {
		return 0, err
	}
	stmt2, err := db.Prepare(`INSERT INTO users (wxUserID) VALUES (?)`)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	defer stmt2.Close()
	res, err := stmt2.Exec(userID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	return int(ID), nil
}

// NewTgUser 添加新的Telegram用户
func NewTgUser(tID int) (int64, error) {
	stmt, err := db.Prepare(`INSERT  INTO tgUsers (id) SELECT (?) FROM DUAL WHERE NOT EXISTS (SELECT id FROM tgUsers WHERE id= ? )`)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tID, tID)
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err
	}
	return int64(tID), nil
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
func BindTg(openID string, tID int) error {
	stmt, err := db.Prepare(`UPDATE users SET tgUserID = ? WHERE wxUserID =（SELECT id FROM wxUsers WHERE openID = ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tID, openID)
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
func UnBindTgFromWX(openID string) error {
	stmt, err := db.Prepare(`UPDATE users SET tgUserID = NULL WHERE wxUserId = (SELECT id FROM wxUsers WHERE openID = ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(openID)
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
func IsBindTg(openid string) (int, error) {
	var tID int
	err := db.QueryRow("SELECT tgUserID FROM users WHERE wxUserID = (SELECT id FROM wxUsers WHERE openID = ?)", openid).Scan(&tID)
	if err == sql.ErrNoRows {
		glog.V(2).Info(fmt.Sprintf("微信openID:%s 用户尚未绑定TG帐号", openid))
		return 0, nil //这个时候是没有绑定
	}
	if err != nil {
		glog.Error("数据库错误：", err)
		return 0, err //出现错误
	}
	return tID, nil
}

// GetUserIDByWx 从wx后去用户ID，如果用户不存在，则新建用户
func GetUserIDByWx(openID, nickName string) int {
	var uID int
	SQL := "SELECT id FROM users WHERE wxUserID=(SELECT id FORM wxUsers WHERE openID = ? LIMIT 1)"
	err := db.QueryRow(SQL, openID).Scan(&uID)
	if err == sql.ErrNoRows {
		w := WxUser{
			openID:   openID,
			NickName: nickName,
		}
		id, err := NewUserFromWx(&w)
		if err != nil {
			return 0
		}
		return id
	} else {
		return uID
	}
}
