package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Id       uint   `json:"id"`
	jwt.StandardClaims
}

// TokenExpireDuration 设置TokenExpireDuration token过期时间为2小时
const TokenExpireDuration = time.Hour * 2

// MySecret 设置一个MySecret 超时以后
var MySecret = []byte("夏天夏天悄悄过去....今年太冷了......")

// GenToken 生成JWT
func GenToken(username string, id uint, uuid string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		uuid,
		username, // 自定义字段
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "ServerManagerSystem",                      // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT（验证token）
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		_, err = Get([]byte(claims.Uuid))
		if err != nil {
			return nil, err
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		//token, _ := c.Cookie("token")
		authHeader := c.Request.Header.Get("Authorization")
		//fmt.Println(token)
		if authHeader == "" {
			resultBody := ResultBody{
				Code:    2003,
				Data:    nil,
				Message: "Token is empty",
			}
			c.JSON(http.StatusNonAuthoritativeInfo, resultBody)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			resultBody := ResultBody{
				Code:    2004,
				Data:    nil,
				Message: "验证方法错误",
			}
			c.JSON(http.StatusNonAuthoritativeInfo, resultBody)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			resultBody := ResultBody{
				Code:    2005,
				Data:    nil,
				Message: "登录过期",
			}
			c.JSON(http.StatusForbidden, resultBody)
			c.Abort()
			return
		}
		_, err = Get([]byte(mc.Uuid))
		if err != nil {
			resultBody := ResultBody{
				Code:    2005,
				Data:    nil,
				Message: "无效登录",
			}
			c.JSON(http.StatusForbidden, resultBody)
			c.Abort()
			return
		}
		//将当前请求的username信息保存到请求的上下文c上
		c.Set("token_username", mc.Username)
		c.Set("token_id", mc.Id)
		c.Set("uuid", mc.Uuid)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// HomeHandler 获取上下文中的当前登录用户信息
func (MyClaims) HomeHandler(c *gin.Context) {
	tokenName := c.MustGet("token_username").(string)
	tokenId := c.MustGet("token_id").(uint)
	uuid := c.MustGet("uuid").(string)
	resultBody := ResultBody{
		Code: 2000,
		Data: gin.H{
			"name": tokenName,
			"id":   tokenId,
			"uuid": uuid,
		},
		Message: "success",
	}
	c.JSON(http.StatusOK, resultBody)
}
