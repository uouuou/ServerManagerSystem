package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"strconv"
	"time"
)

// Decimal 保留两位小数的方法
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// AutoInfo 计算auth参数
func AutoInfo(token string, userid string) bool {
	timeNow := time.Now().Format("2006-01-02 15")
	tokenNow := Md5V(userid + "ServerManagerSystem2021" + timeNow)
	return token == tokenNow
}

// Md5V 生成32位md5字串
func Md5V(str string) string {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		mid.Log().Error(fmt.Sprintf("err:%v", err))
	}
	return hex.EncodeToString(h.Sum(nil))

}
