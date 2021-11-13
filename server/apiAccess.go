package server

import (
	"context"
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"net/http"
	"os"
	"time"
)

// GetApi 一个Git接口访问方法
func GetApi(url string, body io.Reader) []byte {
	var resultBody []byte
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		mid.Log.Info(fmt.Sprintf("err:%v", err))
	}
	cancel, res, err := fetch(req)
	if err != nil {
		mid.Log.Info(fmt.Sprintf("err:%v", err))
		return nil
	}
	defer cancel()
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode == http.StatusOK {
		resultBody, err = io.ReadAll(res.Body)
		if err != nil {
			mid.Log.Info(fmt.Sprintf("err:%v", err))
		}
	} else {
		return nil
	}
	return resultBody
}

// Down 下载文件到对应位置
func Down(url string, path string) bool {
	resp, err := http.Get(url)
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + ":下载异常 " + err.Error())
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 创建一个文件用于保存
	out, err := os.Create(path)
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + ":保存异常 " + err.Error())
		return false
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + ":文件写入异常 " + err.Error())
		return false
	}
	return true
}

//一个给予http的访问控制方法，用于控制超时时间
func fetch(req *http.Request) (context.CancelFunc, *http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	resp, err := ctxhttp.Do(ctx, http.DefaultClient, req)
	return cancel, resp, err
}
