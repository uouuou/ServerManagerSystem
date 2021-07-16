package middleware

import (
	"log"

	loggers "github.com/phachon/go-logger"
)

func Log() *loggers.Logger {
	logger := loggers.NewLogger()

	err := logger.Detach("console")
	if err != nil {
		log.Print(err)
	}
	// 命令行输出配置
	consoleConfig := &loggers.ConsoleConfig{
		Color:      true,  // 命令行输出字符串是否显示颜色
		JsonFormat: false, // 命令行输出字符串是否格式化
		Format:     "",    // 如果输出的不是 json 字符串，JsonFormat: false, 自定义输出的格式
	}
	// 添加 console 为 logger 的一个输出
	err = logger.Attach("console", loggers.LOGGER_LEVEL_DEBUG, consoleConfig)
	if err != nil {
		log.Print(err)
	}
	// 文件输出配置
	fileConfig := &loggers.FileConfig{
		Filename: Dir + "/log/console.log", // 日志输出文件名，不自动存在
		// 如果要将单独的日志分离为文件，请配置LevelFileNem参数。
		LevelFileName: map[int]string{
			logger.LoggerLevel("error"):   Dir + "/log/console_error.log",   // Error 级别日志被写入 console_error .log 文件
			logger.LoggerLevel("info"):    Dir + "/log/console_info.log",    // Info 级别日志被写入到 console_info.log 文件中
			logger.LoggerLevel("debug"):   Dir + "/log/console_debug.log",   // Debug 级别日志被写入到 console_debug.log 文件中
			logger.LoggerLevel("warning"): Dir + "/log/console_warning.log", // Debug 级别日志被写入到 console_warning.log 文件中
		},
		MaxSize:    1024 * 1024, // 文件最大值（KB），默认值0不限
		MaxLine:    100000,      // 文件最大行数，默认 0 不限制
		DateSlice:  "d",         // 文件根据日期切分， 支持 "Y" (年), "m" (月), "d" (日), "H" (时), 默认 "no"， 不切分
		JsonFormat: true,        // 写入文件的数据是否 json 格式化
		Format:     "",          // 如果写入文件的数据不 json 格式化，自定义日志格式
	}
	// 添加 file 为 logger 的一个输出
	err = logger.Attach("file", loggers.LOGGER_LEVEL_DEBUG, fileConfig)
	if err != nil {
		log.Print(err)
	}
	return logger
}
