package middleware

import (
	"embed"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/hprose/hprose-golang/v3/rpc/plugins/push"
	"github.com/hprose/hprose-golang/v3/rpc/plugins/reverse"
)

var (
	ReadType     string          // sftp读取文本的后缀名使用 | 隔开
	Mode         string          // 运行类型 默认是server服务端 输入-m client就是客户端模式启动
	Version      bool            // 版本号flag输入-v即可触发方法
	Readable     bool            // 检测操作系统的可读性 -r
	flagSet      map[string]bool // flag设置
	CUId         string          // 客户端UUID
	Server       string          // 远端服务端RPC地址
	IDList       []string        // RPC客户端连接列表
	NPSConfig    string          // 客户端NPS配置信息
	FRPConfig    string          // 客户端FRP配置信息
	ServerHttp   string          // 远程服务端HTTP地址
	Auth         string          // 客户端连接远程服务端使用的加密秘药
	Service      *push.Broker    // RPC服务端
	Caller       *reverse.Caller // RPC Caller服务端
	MainName     string          // 主程序名称
	NatAuth      int             // 远程对客户端的NAT功能控制 1为开启 2为关闭
	CronAuth     int             // 远程对客户端Cron的控制 1为开启 2为关闭
	FS           embed.FS        // 文件打包后的路径
	AppRunStatus []AppRunStart
)

// AppRunStart 用于记录process模块启动的数据
type AppRunStart struct {
	Name   string `json:"name"`
	Status bool   `json:"AppRunStart"`
	Num    int    `json:"num"`
	Msg    bool   `json:"msg_num"`
}

func ReadFileType() string {
	return ReadType
}

func GetMode() string {
	return Mode
}

func GetAuth() string {
	return Auth
}

func GetCUId() string {
	return CUId
}

func GetServer() string {
	return Server
}

func GetIDList() []string {
	IDList = GetRpcService().IdList("test")
	return IDList
}

func GetNPSConfig() string {
	return NPSConfig
}

func GetFRTPConfig() string {
	return FRPConfig
}

func GetServerHttp() string {
	return ServerHttp
}

func GetRpcService() *push.Broker {
	return Service
}

func GetNatAuth() int {
	return NatAuth
}

func GetCronAuth() int {
	return CronAuth
}

func GetTokenName(c *gin.Context) string {
	name, _ := c.Get("token_username")
	return name.(string)
}

func init() {
	flag.BoolVar(&Version, "v", false, "Print program build version")
	flag.BoolVar(&Readable, "r", false, "Readable when testing the system environment")
	flag.StringVar(&Mode, "m", "server", "Set the boot type, the client starts using the client server using server")
	if !flag.Parsed() {
		flag.Parse()
	}
	flagSet = map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		flagSet[f.Name] = true
	})
}
