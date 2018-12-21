package server

import (
	"strconv"
	"tg_gif/model"

	"github.com/gin-gonic/gin"
)

// 可以单独独立出来的服务，先放这儿吧

// Login 用户登录接口，返回一个生成的jwtToken，但是不包含过期时间?
// 参数jscode，url参数，get
func Login(c *gin.Context) {
	jscode := c.DefaultQuery("jscode", "")
	nickName := c.DefaultQuery("nickName", "")
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
		c.AbortWithStatus(400)
		return
	}
	inttgID, err := strconv.Atoi(tgID)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	if model.BindTg(openID, inttgID) != nil {
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
//参数？index=1&count=10?ask=11
func GetMyGifs(c *gin.Context) {
	index := c.DefaultQuery("index", "1")
	count := c.DefaultQuery("count", "10")
	uid := c.DefaultQuery("ask", "0")
	// openID := c.MustGet("openid").(string)
	iindex, err := strconv.Atoi(index)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	icount, err := strconv.Atoi(count)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	iuid, err := strconv.Atoi(uid)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	gifs := model.GetGifs(iindex, icount, iuid)
	c.JSON(200, gifs)
	return
}