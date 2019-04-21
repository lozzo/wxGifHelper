package server

import (
	"strconv"
	"tg_gif/model"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// 可以单独独立出来的服务，先放这儿吧

// Login 用户登录接口，返回一个生成的jwtToken，但是不包含过期时间?
// 参数jscode，url参数，get
func Login(c *gin.Context) {
	jscode := c.DefaultQuery("jscode", "")
	nickName := c.DefaultQuery("nickName", "")
	glog.V(5).Info(jscode, nickName)
	if jscode == "" {
		c.AbortWithStatus(401)
		return
	}
	openid, err := GetWxOpenID(jscode)
	if err != nil || openid == "" {
		c.AbortWithStatus(401)
		return
	}
	jwtToken := CreatToken(openid)
	tgID, err := model.IsBindTg(openid)
	userID := model.GetUserIDByWx(openid, nickName)
	if err != nil || userID == 0 {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{
		"token":  jwtToken,
		"tgID":   tgID, // 0代表没有绑定
		"userID": userID,
	})
	return
}

// BindTg wx用户绑定Tg帐号
// 参数tgID，url参数，get
func BindTg(c *gin.Context) {
	openID := c.MustGet("openid").(string)
	tgID := c.DefaultQuery("tgID", "")
	if tgID == "" {
		glog.V(5).Info("tgID参数错误")
		c.AbortWithStatus(400)
		return
	}
	inttgID, err := strconv.Atoi(tgID)
	if err != nil {
		glog.V(5).Info("tgID参数类型错误")
		c.AbortWithStatus(500)
		return
	}
	if model.BindTg(openID, inttgID) != nil {
		glog.V(5).Info("绑定失败")
		c.AbortWithStatus(500)
		return
	}
	c.Status(200)
}

// UnBindTg 解除tg绑定，无参数。get
func UnBindTg(c *gin.Context) {
	openID := c.MustGet("openid").(string)
	if model.UnBindTgFromWX(openID) != nil {
		c.AbortWithStatus(500)
		return
	}
	c.Status(200)
	return
}

// GetMyGifs 获取用户的表情包
//参数?ask=11
func GetMyGifs(c *gin.Context) {
	uid := c.DefaultQuery("ask", "0")
	// openID := c.MustGet("openid").(string)

	iuid, err := strconv.Atoi(uid)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	gifs := model.GetGifs(iuid)

	c.JSON(200, gin.H{
		"gifs": gifs,
	})
	return
}

// GetRandGifs 获取随机的表情
func GetRandGifs(c *gin.Context) {
	gifs := model.GetRandGifs(201)
	c.JSON(200, gin.H{
		"gifs": gifs,
	})
	return
}

// ReportGifs 举报有问题的图片
// 参数?id=CAADBQADVQIAAiix4w0FzkQef-eN5QI
func ReportGifs(c *gin.Context) {
	FileID := c.DefaultQuery("id", "0")
	if FileID == "0" {
		c.AbortWithStatus(400)
		return
	}
	model.ReportGifs(FileID)
	c.Status(200)
	return
}

// DeleteUserFile 删除用户文件
// 参数?id=CAADBQADVQIAAiix4w0FzkQef-eN5QI&ask=11
func DeleteUserFile(c *gin.Context) {
	FileID := c.DefaultQuery("id", "0")
	if FileID == "0" {
		c.AbortWithStatus(400)
		return
	}
	uid := c.DefaultQuery("ask", "0")
	if uid == "0" {
		c.AbortWithStatus(400)
		return
	}
	model.DeleteUserFile(FileID, uid)
	c.Status(200)
	return
}

// SetToMyGifs 添加到我喜欢的gifs内
// 参数?id=CAADBQADVQIAAiix4w0FzkQef-eN5QI&ask=11
func SetToMyGifs(c *gin.Context) {
	FileID := c.DefaultQuery("id", "0")
	if FileID == "0" {
		c.AbortWithStatus(400)
		return
	}
	uid := c.DefaultQuery("ask", "0")
	if uid == "0" {
		c.AbortWithStatus(400)
		return
	}
	model.AddFilesFromWx(FileID, uid)
	c.Status(200)
	return
}
