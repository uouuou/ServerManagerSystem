package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// AutoInfo 内部方法的token验证
func AutoInfo(token string, userid string) bool {
	timeNow := time.Now().Format("2006-01-02 15")
	tokenNow := Md5V(userid + "ServerManagerSystem2021" + timeNow)
	return token == tokenNow
}

// Md5V 获取字符串的MD5值
func Md5V(str string) string {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		Log().Error(fmt.Sprintf("err:%v", err))
	}
	return hex.EncodeToString(h.Sum(nil))

}

// CId 用户客户端ID获取方法
func CId() string {
	return GetUUID()
}
