package system

import (
	mod "github.com/uouuou/ServerManagerSystem/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// ResponseBody 结构体
type ResponseBody struct {
	Code     int         `json:"code"`
	Duration string      `json:"duration"`
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
}

type speedInfo struct {
	Up   uint64
	Down uint64
}

var si *speedInfo

// TimeCost web函数执行用时统计方法
func TimeCost(start time.Time, body *ResponseBody) {
	body.Duration = time.Since(start).String()
}

// CollectTask 启动收集主机信息任务
func CollectTask() {
	var recvCount, sentCount uint64
	c := cron.New()
	lastIO, _ := net.IOCounters(true)
	var lastRecvCount, lastSentCount uint64
	for _, k := range lastIO {
		lastRecvCount = lastRecvCount + k.BytesRecv
		lastSentCount = lastSentCount + k.BytesSent
	}
	si = &speedInfo{}
	_, _ = c.AddFunc("@every 2s", func() {
		result, _ := net.IOCounters(true)
		recvCount, sentCount = 0, 0
		for _, k := range result {
			recvCount = recvCount + k.BytesRecv
			sentCount = sentCount + k.BytesSent
		}
		si.Up = (sentCount - lastSentCount) / 2
		si.Down = (recvCount - lastRecvCount) / 2
		lastSentCount = sentCount
		lastRecvCount = recvCount
		lastIO = result
	})
	c.Start()
}

// ServerInfo 获取服务器信息
func ServerInfo() *ResponseBody {
	var nftVersions string
	var unboundVersions string
	var clashVersions string
	responseBody := ResponseBody{Message: "success"}
	defer TimeCost(time.Now(), &responseBody)
	cpuPercent, _ := cpu.Percent(0, false)
	cpuInfo, _ := cpu.Info()
	vmInfo, _ := mem.VirtualMemory()
	smInfo, _ := mem.SwapMemory()
	diskInfo, _ := disk.Usage("/")
	loadInfo, _ := load.Avg()
	tcpCon, _ := net.Connections("tcp")
	udpCon, _ := net.Connections("udp")
	netCount := map[string]int{
		"tcp": len(tcpCon),
		"udp": len(udpCon),
	}
	timestamp, _ := host.BootTime()
	kernelVersion, _ := host.KernelVersion()
	platform, family, systemVersion, _ := host.PlatformInformation()
	t := time.Unix(int64(timestamp), 0)
	bootTime := t.Local().Format("2006-01-02 15:04:05")
	netInfo, _ := net.Interfaces()
	//检测clash版本号
	//检测clash版本号
	infoClash, err := mod.ExecCommandWithResult("clash -v")
	if err == nil {
		clashVersion := strings.Fields(strings.TrimSpace(infoClash))
		clashVersions = clashVersion[1]
	}
	//检测unbound的安装
	infoUnbound, err := mod.ExecCommandWithResult("unbound -v")
	if err == nil {
		unboundVersion := strings.Fields(strings.TrimSpace(infoUnbound))
		unboundVersions = unboundVersion[6]
	}
	//检测nft是否安装并输出版本号
	info, err := mod.ExecCommandWithResult("nft -v")
	if err == nil {
		ntfVersion := strings.Fields(strings.TrimSpace(info))
		nftVersions = ntfVersion[1]
	}
	ioSata, _ := disk.IOCounters()
	responseBody.Data = map[string]interface{}{
		"cpu":            cpuInfo,
		"cpu_used":       cpuPercent,
		"memory":         vmInfo,
		"swap":           smInfo,
		"net":            netInfo,
		"disk":           diskInfo,
		"load":           loadInfo,
		"speed":          si,
		"netCount":       netCount,
		"runtime":        bootTime,
		"kernel_version": kernelVersion,
		"system_version": systemVersion,
		"platform":       platform,
		"family":         family,
		"nftables":       nftVersions,
		"unbound":        unboundVersions,
		"clash":          clashVersions,
		"io":             ioSata,
	}
	responseBody.Code = 2000
	return &responseBody
}

// GetServerInfo 获取服务器信息的接口
func GetServerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, ServerInfo())
}
