package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid := c.Query("userid")
		token := c.Query("token")
		if AutoInfo(token, userid) == true {
			c.Next()
		} else {
			c.Abort()
			c.JSON(200, gin.H{
				"code":    0,
				"status":  false,
				"message": "err",
			})
		}
	}
}

func HandleNotFound(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/404")
}
