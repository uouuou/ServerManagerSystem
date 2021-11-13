package nat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	con "github.com/uouuou/ServerManagerSystem/server"
	"os"
	"strings"
	"time"
)

// NPS配置文件读取状态
var starts bool

// GitHubReleasesApi GitHubReleasesApi需要的结构体
type GitHubReleasesApi struct {
	TagName   string    `json:"tag_name"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Assets    []struct {
		Name               string    `json:"name"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadUrl string    `json:"browser_download_url"`
	} `json:"assets"`
}

// Nps 在线Nps信息
type Nps struct {
	PathName           string   `json:"path_name"`
	TagName            string   `json:"tag_name"`
	Name               string   `json:"name"`
	UpdatedAt          mod.Time `json:"updated_at"`
	Path               string   `json:"path"`
	BrowserDownloadUrl string   `json:"browser_download_url"`
	Starts             string   `json:"starts" yaml:"Nps"`
}

// OnlineNpsVersion 获取线上的Nps版本信息
func OnlineNpsVersion(c *gin.Context) {
	var (
		nps Nps
	)
	nps = GetNpsNew(mid.GetMode())
	mid.DataOk(c, nps, "查询成功")
}

// UpdateNpsNew 更新Nps版本
func UpdateNpsNew(c *gin.Context) {
	var (
		nps Nps
	)
	nps = GetNpsNew(mid.GetMode())
	if UpdateNps(nps.Starts) {
		mid.DataOk(c, nps, "更新成功")
	} else {
		mid.DataNot(c, nps, "暂无更新")
	}
}

// UpdateNps 更新Nps版本
func UpdateNps(start string) bool {
	var (
		nps Nps
	)
	nps = GetNpsNew(start)
	info := mod.NpsVersion()
	if info == "" {
		info = "Info: Nps New Install"
	}
	if nps.BrowserDownloadUrl == "" {
		mid.Log.Info(mid.RunFuncName() + ":Nps接口请求异常")
		return false
	}
	npsInfo := strings.Split(info, ": ")
	if "v"+npsInfo[1] != nps.TagName {
		url := fmt.Sprintf("%v%v", "https://mirror.ghproxy.com/", nps.BrowserDownloadUrl)
		if con.Down(url, nps.Path) {
			mid.Log.Info("Nps " + nps.TagName + " 下载成功：" + nps.PathName)
		}
		pid := mod.ForPid("nps")
		pids := strings.Split(pid, "\n")
		if pids[0] == "" {
			if !mod.IsExists(mid.Dir + "/config/nps") {
				err := os.MkdirAll(mid.Dir+"/config/nps", 0755)
				if err != nil {
					return false
				}
			}
			err := mod.ExecCommand("tar -zvxf " + mid.Dir + "/config/nps.gz -C " + mid.Dir + "/config/nps")
			if err != nil {
				mid.Log.Error(err.Error())
			}
		} else {
			if !mod.IsExists(mid.Dir + "/config/nps") {
				err := os.MkdirAll(mid.Dir+"/config/nps", 0755)
				if err != nil {
					return false
				}
			}
			for _, s := range pids {
				if s != " " {
					_, _ = mod.ExecCommandWithResult(fmt.Sprintf("kill -9 %v", s))
				} else {
					continue
				}
			}
			err := mod.ExecCommand("tar -zvxf " + mid.Dir + "/config/nps.gz -C " + mid.Dir + "/config/nps")
			if err != nil {
				mid.Log.Error(err.Error())
			}
		}
		err := mod.ExecCommand("chmod -R 755 " + mid.Dir + "/config/nps")
		if err != nil {
			mid.Log.Error(err.Error())
		}
		return true
	} else {
		return false
	}
}

// GetNpsNew 获取对应操作系统系统的Nps信息
func GetNpsNew(starts string) (Nps Nps) {
	var (
		gitApi GitHubReleasesApi
		url    string
	)
	url = "https://api.github.com/repos/ehang-io/nps/releases/latest"
	npsInfo := con.GetApi(url, nil)
	err := json.Unmarshal(npsInfo, &gitApi)
	if err != nil {
		mid.Log.Error("NPS线上接口请求受限：" + err.Error())
		return
	}
	architecture, system := mid.GetArchitecture()
	Nps.Path = mid.Dir + "/config/nps.gz"
	Nps.TagName = gitApi.TagName
	Nps.UpdatedAt = mod.Time(gitApi.CreatedAt)
	for _, s := range gitApi.Assets {
		name := strings.Split(s.Name, "_")
		if len(name) >= 3 {
			start := strings.Split(name[2], ".")
			if name[0] == system && name[1] == architecture && start[0] == starts {
				Nps.BrowserDownloadUrl = s.BrowserDownloadUrl
				Nps.PathName = s.Name
			}
		}

	}
	Nps.Starts = starts
	Nps.Name = gitApi.Name
	return
}

// MoNpsConfig 基础设置npc.conf文件并返回
func MoNpsConfig(start string) []byte {
	var name string
	if start == "client" {
		name = "npc.conf"
	}
	if start == "server" {
		name = "/conf/nps.conf"
	}
	config, err := os.ReadFile(mid.Dir + "/config/nps/" + name)
	if err != nil {
		if !starts {
			mid.Log.Error(mid.RunFuncName() + ":未找到配置文件，等待线上发布....")
			starts = true
		}
		return config
	}
	if start == "client" {
		server := strings.Split(mid.Server, "://")
		ip := strings.Split(server[1], ":")
		vkey := strings.Split(string(config), ":")
		config = []byte(ip[0] + ":" + vkey[0] + " -vkey=" + vkey[1])
	}
	return config
}

//SetNpsConfig 设置NPS配置文件
func SetNpsConfig(c *gin.Context) {
	type frp struct {
		Data string `json:"data"`
	}
	var frps frp
	err := c.ShouldBind(&frps)
	if err != nil {
		mid.Log.Error(err.Error())
		mid.ClientErr(c, err, "格式错误")
		return
	}
	err = os.WriteFile(mid.Dir+"/config/nps/conf/nps.conf", []byte(frps.Data), 0755)
	if err != nil {
		mid.DataNot(c, nil, "写入出错")
	} else {
		mid.FRPConfig = frps.Data
		mid.DataOk(c, frps, "正常写入")
	}
}

// GetNpsConfig 获取MPS 安装情况及配置文件
func GetNpsConfig(c *gin.Context) {
	readFile, err := os.ReadFile(mid.Dir + "/config/nps/conf/nps.conf")
	if err != nil {
		mid.DataErr(c, err, "文件读取错误")
		return
	}
	mid.DataOk(c, string(readFile), "查询成功")
}
