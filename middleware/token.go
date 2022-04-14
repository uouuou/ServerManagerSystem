package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/duke-git/lancet/convertor"
	"time"
)

// AutoInfo 内部方法的token验证
func AutoInfo(token string, userid string, timestamp string) bool {
	toInt, err := convertor.ToInt(timestamp)
	if err != nil {
		return false
	}
	if time.Now().Unix()-toInt > 3600 && time.Now().Unix()-toInt < -3600 {
		return false
	}
	tokenNow := Md5V(userid + "ServerManagerSystem2021" + timestamp)
	return token == tokenNow
}

// Md5V 获取字符串的MD5值
func Md5V(str string) string {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		Log.Error(fmt.Sprintf("err:%v", err))
	}
	return hex.EncodeToString(h.Sum(nil))

}

// CId 用户客户端ID获取方法
func CId() string {
	return GetUUID()
}
