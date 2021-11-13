package models

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/middleware/convert"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// PortIsUse 判断端口是否占用
func PortIsUse(port int) bool {
	_, tcpError := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), time.Millisecond*50)
	udpAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	udpConn, udpError := net.ListenUDP("udp", udpAddr)
	if udpConn != nil {
		defer func(udpConn *net.UDPConn) {
			_ = udpConn.Close()
		}(udpConn)
	}
	return tcpError == nil || udpError != nil
}

// RandomPort 获取没占用的随机端口
func RandomPort() int {
	for {
		rand.Seed(time.Now().UnixNano())
		newPort := rand.Intn(65536)
		if !PortIsUse(newPort) {
			return newPort
		}
	}
}

// IsExists 检测指定路径文件或者文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// GetLocalIP 获取本机ipv4地址
func GetLocalIP() string {
	resp, err := http.Get("http://api.ipify.org")
	if err != nil {
		resp, _ = http.Get("http://icanhazip.com")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	s, _ := io.ReadAll(resp.Body)
	return string(s)
}

// CheckIP 检测ipv4地址的合法性
func CheckIP(ip string) bool {
	isOk, err := regexp.Match(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, []byte(ip))
	if err != nil {
		fmt.Println(err)
	}
	return isOk
}

// InstallPack 安装指定名字软件
func InstallPack(name string) {
	if !CheckCommandExists(name) {
		if CheckCommandExists("yum") {
			err := ExecCommand("yum install -y " + name)
			if err != nil {
				mid.Log.Error(err.Error())
			}
		} else if CheckCommandExists("apt-get") {
			err := ExecCommand("apt-get update")
			if err != nil {
				mid.Log.Error(err.Error())
			}
			err = ExecCommand("apt-get install -y " + name)
			if err != nil {
				mid.Log.Error(err.Error())
			}
		}
	}
}

// OpenPort 开通指定端口
func OpenPort(port int) {
	if CheckCommandExists("firewall-cmd") {
		err := ExecCommand(fmt.Sprintf("firewall-cmd --zone=public --add-port=%d/tcp --add-port=%d/udp --permanent >/dev/null 2>&1", port, port))
		if err != nil {
			mid.Log.Error(err.Error())
		}
		err = ExecCommand("firewall-cmd --reload >/dev/null 2>&1")
		if err != nil {
			mid.Log.Error(err.Error())
		}
	} else {
		withResult, err := ExecCommandWithResult(fmt.Sprintf(`iptables -nvL --line-number|grep -w "%d"`, port))
		if err != nil {
			mid.Log.Error(err.Error())
		}
		if len(withResult) > 0 {
			return
		}
		err = ExecCommand(fmt.Sprintf("iptables -I INPUT -p tcp --dport %d -j ACCEPT", port))
		if err != nil {
			mid.Log.Error(err.Error())
		}
		err = ExecCommand(fmt.Sprintf("iptables -I INPUT -p udp --dport %d -j ACCEPT", port))
		if err != nil {
			mid.Log.Error(err.Error())
		}
		err = ExecCommand(fmt.Sprintf("iptables -I OUTPUT -p udp --sport %d -j ACCEPT", port))
		if err != nil {
			mid.Log.Error(err.Error())
		}
		err = ExecCommand(fmt.Sprintf("iptables -I OUTPUT -p tcp --sport %d -j ACCEPT", port))
	}
}

// ForPid 获取linux应用程序运行id
func ForPid(name string) (pid string) {
	pid, err := ExecCommandWithResult(fmt.Sprintf("ps aux|grep '%v'|grep -v \"grep\"|awk '{print $2}'", name))
	if err != nil {
		mid.Log.Error("获取PID异常")
		return
	}
	pid = strings.Replace(pid, "\n", "", -1)
	return
}

// ForPids 获取linux应用程序运行id(可能有多个)
func ForPids(name string) (pids []string) {
	pid, err := ExecCommandWithResult(fmt.Sprintf("ps aux|grep '%v'|grep -v \"grep\"|awk '{print $2}'", name))
	if err != nil {
		mid.Log.Error("获取PID异常")
		return
	}
	pid = strings.Replace(pid, "\n", ",", -1)
	pids = strings.Split(pid, ",")
	return
}

// ForPidString 获取linux应用程序运行id(可能有多个)
func ForPidString(name string) (pid string) {
	pid, err := ExecCommandWithResult(fmt.Sprintf("ps aux|grep '%v'|grep -v \"grep\"|awk '{print $2}'", name))
	if err != nil {
		mid.Log.Error("获取PID异常")
		return
	}
	pid = strings.Replace(pid, "\n", ",", -1)
	pid = convert.TrimLastChar(pid)
	return
}

// CheckProRunning 根据进程名判断进程是否运行
func CheckProRunning(serverName string) (bool, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	if err != nil {
		return false, err
	}
	return pid != "", nil
}

// GetPid 根据进程名称获取进程ID
func GetPid(serverName string) (string, error) {
	a := `ps ux | awk '/` + serverName + `/ && !/awk/ {print $2}'`
	pid, err := RunCommand(a)
	return pid, err
}

// 在win上运行程序
func runInWindows(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

// RunCommand 检测操作系统决定在哪里运行cmd
func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}

//在linux上运行cmd
func runInLinux(cmd string) (string, error) {
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

// GetCurrentDirectory 程序运行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		mid.Log.Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// ClashVersion 检测clash版本号
func ClashVersion() string {
	var ClashVersion string
	infoClash, err := ExecCommandWithResult("clash -v")
	if err != nil {
		ClashVersion = ""
	} else {
		clashVersion := strings.Fields(strings.TrimSpace(infoClash))
		ClashVersion = fmt.Sprintf("%v %v", clashVersion[0], clashVersion[1])
	}
	return ClashVersion
}

// NpsVersion 检测nps版本号
func NpsVersion() string {
	var (
		NpsVersion string
		name       string
	)
	switch mid.GetMode() {
	case "client":
		name = "npc"
	case "server":
		name = "nps"
	}
	info, err := ExecCommandWithResult(mid.Dir + "/config/nps/" + name + " -version")
	if err != nil {
		NpsVersion = ""
	} else {
		clashVersion := strings.Fields(strings.TrimSpace(info))
		NpsVersion = fmt.Sprintf("%v %v", clashVersion[0], clashVersion[1])
	}
	return NpsVersion
}

// FrpVersion 检测Frp版本号
func FrpVersion() string {
	var (
		FrpVersion string
		name       string
	)
	switch mid.GetMode() {
	case "client":
		name = "frpc"
	case "server":
		name = "frps"
	}
	info, err := ExecCommandWithResult(mid.Dir + "/config/frp/" + name + " -v")
	if err != nil {

		FrpVersion = ""
	} else {
		infos := strings.Split(info, "\n")
		FrpVersion = infos[0]
	}
	return FrpVersion
}

// FileExist 获取某个文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
