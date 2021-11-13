package middleware

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	Log *logrus.Logger
)

// InitLogger 系统日志准备
func InitLogger() {
	Log = logrus.New()
	//创建Hook
	hook := NewLfsHook(filepath.Join("log", "ServerManagerSystem"), 24*time.Hour, 15)
	Log.AddHook(hook)
	//日志格式化
	Log.SetFormatter(formatter(true))
	//日志级别
	Log.SetLevel(logrus.InfoLevel)
	//是否开启显示方法
	Log.SetReportCaller(true)
}

// GinLog Gin的日志记录
func GinLog() gin.HandlerFunc {
	logger := logrus.New()
	//日志格式化
	logger.SetFormatter(formatter(true))
	//日志级别
	logger.SetLevel(logrus.InfoLevel)
	//是否开启显示方法
	logger.SetReportCaller(false)
	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Log.Error(err.Error())
	}
	logger.Out = src
	//创建Hook
	hook := NewLfsHook(filepath.Join("log", "ServerManagerSystemGinLog"), 24*time.Hour, 15)
	logger.AddHook(hook)
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//hostname
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求长度
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		// 请求IP
		clientIP := c.ClientIP()
		//用户浏览器
		userAgent := c.Request.UserAgent()
		// 日志格式
		entry := logger.WithFields(logrus.Fields{
			"HostName":    hostName,
			"Status":      statusCode,
			"LatencyTime": latencyTime,
			"Ip":          clientIP,
			"Method":      reqMethod,
			"Path":        reqUri,
			"DataSize":    dataSize,
			"UserAgent":   userAgent,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if len(c.Errors) >= 500 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else if statusCode >= 400 {
			entry.Warn(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			entry.Info(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}

//日志格式化设置
func formatter(isConsole bool) *nested.Formatter {
	formatTer := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	}
	if isConsole {
		formatTer.NoColors = false
	} else {
		formatTer.NoColors = true
	}
	return formatTer
}

//NewLfsHook 位置Hook设置
func NewLfsHook(logName string, rotationTime time.Duration, leastDay uint) logrus.Hook {
	writer, err := rotatelogs.New(
		// 日志文件
		logName+".%Y%m%d%H%M%S"+".log",
		// 日志周期(默认每86400秒/一天旋转一次)
		rotatelogs.WithRotationTime(rotationTime),
		// 生成软链，指向最新日志文件
		//rotatelogs.WithLinkName(logName+".log"),
		// 清除历史 (WithMaxAge和WithRotationCount只能选其一)
		//rotatelogs.WithMaxAge(time.Hour*24*7), //默认每7天清除下日志文件
		rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
	)
	if err != nil {
		panic(err)
	}
	// 可设置按不同level创建不同的文件名
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
		//}, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	}, formatter(false))

	return lfsHook
}
