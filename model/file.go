package model

import (
	"tg_gif/common"

	"github.com/golang/glog"
)

// 表情包相关的服务

// AddFiles 批量添加gifs文件，通用的添加文件
func AddFiles(gifs []*Gifs) error {
	tx, err := db.Begin()
	if err != nil {
		glog.Error("数据库错误：", err)
		return err
	}
	isGroup := false
	groupID := 0
	if gifs[0].Group != nil {
		groupID, err = addGroup(gifs[0].Group.Name)
		if err != nil {
			return err
		}
	}
	for _, gif := range gifs {
		if isGroup {
			sql := "INSERT INTO gifs (GroupID,FileID,UserID) VALUES (?,?,?)"
			_, err := tx.Exec(sql, groupID, gif.File, gif.User.ID)
			if err != nil {
				tx.Rollback()
				return err

			}
		} else {
			sql := "INSERT INTO gifs (FileID,UserID) VALUES (?,?)"
			_, err := tx.Exec(sql, gif.File, gif.User.ID)
			if err != nil {
				tx.Rollback()
				return err

			}
		}
	}
	err = tx.Commit()
	if err != nil {
		panic("数据库错误")
	}
	return nil
}

//AddFilesFromTg 将MsgStatus内的表情写入数据库
func AddFilesFromTg(m *common.MsgStatus) {
	tx, err := db.Begin()
	if err != nil {
		glog.Error("数据库错误：", err)
	}
	for _, file := range m.File {
		sql := "INSERT INTO gifs (FileID,UserID) VALUES (?,( SELECT id FROM users WHERE tgUserID = ? LIMIT 1)) "
		// sql := "INSERT INTO gifs (FileID,UserID) VALUES (?,1)"
		_, err := tx.Exec(sql, file.ID, m.ID)
		if err != nil {
			tx.Rollback()
			glog.V(5).Info(file.ID, m.ID)
			glog.Error("数据库错误：", err)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		panic("数据库错误")
	}
}

func addGroup(name string) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO gifGroups (name) VALUES (?)`)
	if err != nil {
		glog.Error("数据库错误:", err)
		return 0, err
	}
	res, err := stmt.Exec(name)
	if err != nil {
		glog.Error("数据库错误:", err)
		return 0, err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		glog.Error("数据库错误:", err)
		return 0, err
	}
	return int(ID), nil
}

// GetGifs 获取一了列表的gifs
func GetGifs(uID int) []string {
	var gifs []string
	SQL := "SELECT FileID FROM gifs WHERE userID = ? GROUP BY FileID ORDER BY max(id) desc LIMIT 200"

	rows, err := db.Query(SQL, uID)
	if err != nil {
		glog.Warning(err)
		return gifs
	}
	defer rows.Close()
	for rows.Next() {
		var gif string
		if err := rows.Scan(&gif); err != nil {
			glog.Warning(err)
			continue
		} else {
			gifs = append(gifs, gif)
		}
	}
	return gifs
}

// GetRandGifs 随机获取一些图片
func GetRandGifs(n int) []string {
	var gifs []string
	SQL := "SELECT DISTINCT FileID FROM gifs ORDER BY rand() LIMIT ?"
	rows, err := db.Query(SQL, n)
	if err != nil {
		glog.Warning(err)
		return gifs
	}
	defer rows.Close()
	for rows.Next() {
		var gif string
		if err := rows.Scan(&gif); err != nil {
			glog.Warning(err)
			continue
		} else {
			gifs = append(gifs, gif)
		}
	}
	return gifs
}

// GetAllFilesID 返回所有非重复的fileID
func GetAllFilesID() []string {
	var gifs []string
	SQL := "SELECT DISTINCT FileID FROM gifs"
	rows, err := db.Query(SQL)
	if err != nil {
		glog.Warning(err)
		return gifs
	}
	defer rows.Close()
	for rows.Next() {
		var gif string
		if err := rows.Scan(&gif); err != nil {
			glog.Warning(err)
			continue
		} else {
			gifs = append(gifs, gif)
		}
	}
	return gifs
}
