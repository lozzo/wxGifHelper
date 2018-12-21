package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

//wx服务端的相关服务

var (
	jwtExp       int64 = 1
	jwtSectetKey       = "^s123(A23S1D41O@I23J1D4UdASDNA(@#jbp"
	appid        string
	secret       string
)

// WxInit 引入appid和secret
func WxInit(appid, secret string) {
	appid = appid
	secret = secret
}

//CreatToken 生成jwttoken
func CreatToken(user string) string {
	stamp := time.Now().Unix()
	mapClaims := jwt.MapClaims{
		"iss": "tg.onemoresec.com", //网站域名
		"iat": stamp,               //签发时间
		"exp": stamp + jwtExp*3600, //过期时间 float64 最后....
		"aud": user,                //签发用户---openID
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte(jwtSectetKey))
	if err != nil {
		return ""
	}
	return tokenString
}

// VerifiJwtToken 验证jwt_token 并返回该token的openid名称
func VerifiJwtToken(tonkenString string) (string, bool) {
	var user string
	myToken, err := jwt.Parse(tonkenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSectetKey), nil
	})
	if err != nil {
		glog.Info("JWTAuth :Parse faild ", err)
		return "", false
	}

	x, ok := myToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	iss := x["iss"]
	exp := x["exp"]
	aud := x["aud"]
	// 校验网站
	if iss != "tg.onemoresec.com" {
		glog.Info("JWTAuth :iss key failed : ", iss)
		return "", false
	}
	// 校验过期事件
	ex, ok := exp.(float64)
	if !ok {
		glog.Info("JWTAuth :exp key failed : ", ok)
		return "", false
	}
	if ex < float64(time.Now().Unix()) {
		glog.Info("JWTAuth :exp Expired")
		return "", false
	}

	//返回user
	au, ok := aud.(string)
	if !ok {
		glog.Info("JWTAuth :aud key failed : ", ok)
		return "", false
	}
	user = au

	if myToken.Valid {
		return user, true
	}
	glog.Info("JWTAuth : not Valid")
	return "", false

}

type getOpenID struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
}

// GetWxOpenID 获取微信openid
func GetWxOpenID(jscode string) (string, error) {
	URL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, jscode)
	openid := getOpenID{}
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &openid)
	if err != nil {
		return "", err
	}
	if openid.OpenID == "" {
		return "", errors.New("blank opeid")
	}
	return openid.OpenID, nil
}
