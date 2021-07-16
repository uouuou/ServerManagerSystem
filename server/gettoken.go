package server

import (
	mod "github.com/uouuou/ServerManagerSystem/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const adminUser = "sms"
const adminPassword = "sms123456"

// GetToken 使用默认的用户名密码采用MD5加密后传入并返回中间件需要的token值
func GetToken(c *gin.Context) {
	user := c.Query("user")
	password := c.Query("password")
	timeNow := time.Now().Format("2006-01-02 15")
	tokenNow := mod.Md5V(user + "ServerManagerSystem2021" + timeNow)
	if user != "" && password != "" {
		if user == adminUser {
			if adminPassword == password {
				c.JSON(http.StatusOK, gin.H{
					"code":    1,
					"status":  true,
					"token":   tokenNow,
					"message": "验证成功！",
				})
				return
			} else {
				c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
					"code":    0,
					"status":  false,
					"token":   "",
					"message": "密码错误！",
				})
				return
			}
		} else {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"code":    0,
				"status":  false,
				"token":   "",
				"message": "用户名不正确！",
			})
			return
		}

	} else {
		c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
			"code":    0,
			"status":  false,
			"token":   "",
			"message": "用户名和密码不能为空！",
		})
	}
}
