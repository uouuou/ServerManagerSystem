package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var Dir string

func init() {
	// 从本地读取环境变量+
	err := godotenv.Load()
	if err != nil {
		gin.SetMode(gin.ReleaseMode)
	}
	if gin.Mode() == gin.ReleaseMode {
		Dir = GetCurrentDirectory()
		MainName = GetMainName()
	} else {
		Dir = "."
		MainName = GetMainName()
	}
	//创建日志文件夹
	if !FileExist(Dir + "/log") {
		if err := os.MkdirAll(Dir+"/log", 0755); err != nil {
			Log.Error(RunFuncName() + ":创建文件夹异常 " + err.Error())
		}
	}
	//创建配置文件夹
	if !FileExist(Dir + "/config") {
		if err := os.MkdirAll(Dir+"/config", 0755); err != nil {
			Log.Error(RunFuncName() + ":创建文件夹异常 " + err.Error())
		}
	}
}

// GetCurrentDirectory 程序运行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log.Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// GetMainName 获取到当前程序名称
func GetMainName() string {
	path, err := os.Executable()
	if err != nil {
		Log.Error(RunFuncName() + "：获取目录出现错误 " + err.Error())
	}
	fileName := filepath.Base(path)
	return fileName
}

// NewRoutine 采用提前recover的方式终止因为goroutine错误导致的整体崩溃
func NewRoutine(f func()) {
	go func() {
		defer func() {
			// Recover from panic.
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Println(err)
				log.Println(stack)
			}
		}()

		f()
	}()
}
