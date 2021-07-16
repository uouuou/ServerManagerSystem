package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ClientBreak code 4001反馈err到接口
func ClientBreak(c *gin.Context, err error, msg string) {
	var data interface{}
	if err == nil {
		data = nil
	} else {
		data = err.Error()
	}
	resultBody := ResultBody{
		Code:    4001,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// ClientErr code 4002反馈err到接口
func ClientErr(c *gin.Context, err error, msg string) {
	var data interface{}
	if err == nil {
		data = nil
	} else {
		data = err.Error()
	}
	resultBody := ResultBody{
		Code:    4002,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// DataErr code 2003反馈err到接口
func DataErr(c *gin.Context, err error, msg string) {
	var data interface{}
	if err == nil {
		data = nil
	} else {
		data = err.Error()
	}
	resultBody := ResultBody{
		Code:    2003,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// DataNot code 2003反馈数据到接口（不是成功数据）
func DataNot(c *gin.Context, data interface{}, msg string) {
	resultBody := ResultBody{
		Code:    2003,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// DataOk code 2000反馈data到接口
func DataOk(c *gin.Context, data interface{}, msg string) {
	resultBody := ResultBody{
		Code:    2000,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// DataInfo 传入code来实现接口
func DataInfo(c *gin.Context, code int, data interface{}, msg string) {
	resultBody := ResultBody{
		Code:    code,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// DataPageOk code 2000反馈大屏page的data到接口
func DataPageOk(c *gin.Context, pages Pages, data interface{}, msg string) {
	resultBody := ResultPageBody{
		Code:    2000,
		Pages:   pages,
		Data:    data,
		Message: msg,
	}
	c.JSON(http.StatusOK, resultBody)
}

// LoginData 登录成功
func LoginData(c *gin.Context, name string, id uint, avatar string, uuid string, tokenString string) {
	resultBody := ResultTokenBody{
		Code: 2000,
		Data: gin.H{
			"name":   name,
			"id":     id,
			"avatar": avatar,
			"uuid":   uuid,
		},
		Token:   tokenString,
		Message: "登录成功",
	}
	c.JSON(http.StatusOK, resultBody)
}

// RpcClientBreak code 4001反馈err到接口
func RpcClientBreak(err error, msg string) (res ResultBody) {
	var data interface{}
	if err == nil {
		data = nil
	} else {
		data = err.Error()
	}
	resultBody := ResultBody{
		Code:    4001,
		Data:    data,
		Message: msg,
	}
	return resultBody
}

// RpcClientErr code 4002反馈err到接口
func RpcClientErr(err error, msg string) (res ResultBody) {
	var data interface{}
	if err == nil {
		data = nil
	} else {
		data = err.Error()
	}
	resultBody := ResultBody{
		Code:    4002,
		Data:    data,
		Message: msg,
	}
	return resultBody
}

// RpcDataErr code 2003反馈err到接口
func RpcDataErr(err error, msg string) (res ResultBody) {
	var data interface{}
	if err == nil {
		data = nil
	} else {
		data = err.Error()
	}
	resultBody := ResultBody{
		Code:    2003,
		Data:    data,
		Message: msg,
	}
	return resultBody
}

// RpcDataNot code 2003反馈数据到接口（不是成功数据）
func RpcDataNot(data interface{}, msg string) (res ResultBody) {
	resultBody := ResultBody{
		Code:    2003,
		Data:    data,
		Message: msg,
	}
	return resultBody
}

// RpcDataOk code 2000反馈data到接口
func RpcDataOk(data interface{}, msg string) (res ResultBody) {
	resultBody := ResultBody{
		Code:    2000,
		Data:    data,
		Message: msg,
	}
	return resultBody
}

// RpcDataOkUp code 2001反馈data到接口
func RpcDataOkUp(data interface{}, msg string) (res ResultBody) {
	resultBody := ResultBody{
		Code:    2001,
		Data:    data,
		Message: msg,
	}
	return resultBody
}

// ResErr 远程调度客户端时的错误返回
func ResErr(err error, msg string) map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["data"] = err.Error()
	m["message"] = msg
	m["code"] = 2003
	return m
}

// ResOk 远程调度客户端时的正确返回
func ResOk(data interface{}, msg string) map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["data"] = data
	m["message"] = msg
	m["code"] = 2000
	return m
}

// ResNotOk 远程调度客户端时的数据异常返回
func ResNotOk(data interface{}, msg string) map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["data"] = data
	m["message"] = msg
	m["code"] = 2003
	return m
}
