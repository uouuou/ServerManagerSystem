package cLogin

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"net/http"
	"time"
)

const adminUser = "sms"
const adminPassword = "sms2021"
const token = "sms2021"

// GetToken 使用默认的用户名密码采用MD5加密后传入并返回中间件需要的token值
func GetToken(c *gin.Context) {
	user := c.Query("user")
	password := c.Query("password")
	timeNow := time.Now().Format("2006-01-02 15")
	tokenNow := mid.Md5V(user + token + timeNow)
	if user != "" && password != "" {
		if user == adminUser {
			pass := mid.Md5V(adminPassword)
			if pass == password {
				c.JSON(http.StatusOK, gin.H{
					"code":   1,
					"status": true,
					"token":  tokenNow,
					"msg":    "验证成功！",
				})
				return
			} else {
				c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
					"code":   0,
					"status": false,
					"token":  "",
					"msg":    "密码错误！",
				})
				return
			}
		} else {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"code":   0,
				"status": false,
				"token":  "",
				"msg":    "用户名不正确！",
			})
			return
		}

	} else {
		c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
			"code":   0,
			"status": false,
			"token":  "",
			"msg":    "用户名和密码不能为空！",
		})
	}
}
