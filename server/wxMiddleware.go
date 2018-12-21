package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

//JWTAuth jwt验证,openid
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			glog.Info("JWTAuth : get no Authorization Token")
			c.AbortWithStatus(401)
			return
		}
		openid, ok := VerifiJwtToken(authorization)
		if ok {
			c.Set("openid", openid)
			c.Next()
		} else {
			c.AbortWithStatus(401)
			return
		}
	}
}
