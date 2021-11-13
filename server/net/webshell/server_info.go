package webshell

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"golang.org/x/crypto/ssh"
	"strings"
)

type HardwareInfo struct {
	Disk    string `json:"hi_disk"`
	Mem     string `json:"hi_mem"`
	NetCard string `json:"hi_net_card"`
	Cpu     string `json:"hi_cpu"`
	System  string `json:"hi_system"`
	Login   string `json:"hi_login"`
	Ps      string `json:"hi_ps"`
	Port    string `json:"hi_port"`
}

// ServerSystemInfo 一个查询对应接入的服务器的运行状态的接口
func ServerSystemInfo(c *gin.Context) {
	var req Request
	err := c.ShouldBind(&req)
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	req.Token = parts[1]
	if err != nil {
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	info, err := CreateHardwareInfo(req.ServerID, 100, 50)
	if err != nil {
		mid.ClientBreak(c, err, "获取信息错误")
		return
	}
	mid.DataOk(c, info, "查询成功")
}

// CreateHardwareInfo 获取对应服务器的运行信息
func CreateHardwareInfo(sid string, cols int, rows int) (hi *HardwareInfo, err error) {
	hi = &HardwareInfo{}
	var (
		loginInfo SshLoginModel
		client    *ssh.Client
	)
	loginInfo = getServerInfo(sid, cols, rows)
	if loginInfo.Addr == "" {
		return
	}
	client, _, _, err = sshConnect(loginInfo)
	if err != nil {
		return
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)
	hi.Disk, err = SshRemoteRunCommand(client, "df -h")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取硬盘失败")
	}
	hi.Mem, err = SshRemoteRunCommand(client, "free -m")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取内存失败")
	}
	hi.NetCard, err = SshRemoteRunCommand(client, "ifconfig")
	if err != nil {
		_, err := SshRemoteRunCommand(client, "apt install net-tools")
		if err == nil {
			hi.NetCard, err = SshRemoteRunCommand(client, "ifconfig")
			if err != nil {
				mid.Log.Error("Server:" + sid + "获取网卡失败")
			}
		} else {
			mid.Log.Error("Server:" + sid + "按钻NET工具失败")
		}
	}
	hi.Cpu, err = SshRemoteRunCommand(client, "cat /proc/cpuinfo")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取CPU失败")
	}

	hi.System, err = SshRemoteRunCommand(client, "uname -a;who -a;")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取系统失败")
	}
	hi.Login, err = SshRemoteRunCommand(client, "w;last")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取Login失败")
	}
	hi.Ps, err = SshRemoteRunCommand(client, "ps -aux")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取ps失败")
	}
	hi.Port, err = SshRemoteRunCommand(client, "netstat -lntp")
	if err != nil {
		mid.Log.Error("Server:" + sid + "获取netstat失败")
	}
	return hi, nil
}

// SshRemoteRunCommand 通过链接ssh远程执行命令
func SshRemoteRunCommand(sshClient *ssh.Client, command string) (string, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	logString := buf.String()
	if err != nil {
		return logString, fmt.Errorf("CMD: %s  OUT: %s  ERROR: %s", command, logString, err)
	}
	return logString, err
}
