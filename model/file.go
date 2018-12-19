package model

import "github.com/golang/glog"

// 表情包相关的服务

// AddFiles 批量添加gifs文件
func AddFiles(gifs []Gifs) error {
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

func addGroup(name string) (int, error) {
	stmt, err := db.Prepare(`INSERT INTO tgUsers (id) VALUES (?)`)
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
