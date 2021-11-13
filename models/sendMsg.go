package models

import (
	"encoding/json"
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/gomail.v2"
)

type result struct {
	Errno   int    `json:"errno"`
	ErrMsg  string `json:"errmsg"`
	Dataset string `json:"dataset"`
}

// SendMail 有错误的时候发送邮件到邮箱提醒
func SendMail(msg string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "TestMail<TestMail@TestMail.com>") //发件人
	m.SetHeader("To", "TestMail@TestMail.com")             //收件人
	m.SetHeader("Subject", "TestMailMsg")                  //邮件标题
	m.SetBody("text/html", msg)                            //邮件内容

	d := gomail.NewDialer("smtp.TestMail.com", 25, "TestMail@TestMail.com", "TestMail2021")
	//邮件发送服务器信息,使用授权码而非密码
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
	//SendServerJan(msg)
}

// SendServerJan 使用server酱的接口发送推送消息到微信
func SendServerJan(msg string) {
	text := "TestMailMsg~~"
	text = url.QueryEscape(text)
	urls := fmt.Sprintf("https://sc.ftqq.com/yourkey.send?text=%s&desp=%s", text, url.QueryEscape(msg))
	res, err := http.Post(urls, "application/json; charset=UTF-8", strings.NewReader(""))
	if nil != err {
		mid.Log.Error(fmt.Sprintf("http post err:%v", err))
		SendMail(err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	// return status
	if http.StatusOK != res.StatusCode {
		msg := fmt.Sprintf("WebService SoapOa request fail, status: %d\n\n", res.StatusCode)
		mid.Log.Warning(msg)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		msg := fmt.Sprintf("Io ReadAll err:%v\n", err)
		mid.Log.Warning(msg)
		return
	}
	var result result
	err = json.Unmarshal(data, &result)
	if err != nil {
		msg := fmt.Sprintf("err:%v\n", err)
		mid.Log.Warning(msg)
	}
	log.Print(result.ErrMsg)
}
