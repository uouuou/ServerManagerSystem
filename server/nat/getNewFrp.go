package nat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	con "github.com/uouuou/ServerManagerSystem/server"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

// Frp 在线FRP信息
type Frp struct {
	PathName           string   `json:"path_name"`
	TagName            string   `json:"tag_name"`
	Name               string   `json:"name"`
	UpdatedAt          mod.Time `json:"updated_at"`
	Path               string   `json:"path"`
	BrowserDownloadUrl string   `json:"browser_download_url"`
	Starts             string   `json:"starts" yaml:"Nps"`
}

// OnlineFrpVersion 获取线上的FRP版本信息
func OnlineFrpVersion(c *gin.Context) {
	var (
		nps Frp
	)
	nps = GetFrpNew(mid.GetMode())
	mid.DataOk(c, nps, "查询成功")
}

// UpdateFrpNew 更新FRP版本
func UpdateFrpNew(c *gin.Context) {
	var (
		frp Frp
	)
	frp = GetFrpNew(mid.GetMode())
	if UpdateFrp(frp.Starts) {
		mid.DataOk(c, frp, "更新成功")
	} else {
		mid.DataNot(c, frp, "暂无更新")
	}
}

//SetFrpConfig 设置FRPS配置文件
func SetFrpConfig(c *gin.Context) {
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
	err = os.WriteFile(mid.Dir+"/config/frp/frps.ini", []byte(frps.Data), 0755)
	if err != nil {
		mid.DataNot(c, nil, "写入出错")
	} else {
		mid.FRPConfig = frps.Data
		mid.DataOk(c, frps, "正常写入")
	}
}

// GetFrpConfig 获取FRPS 安装情况及配置文件
func GetFrpConfig(c *gin.Context) {
	config := FrpConfig(mid.GetMode())
	mid.DataOk(c, string(config), "查询成功")
}

// FrpConfig 基础设置frps.ini文件并返回
func FrpConfig(start string) []byte {
	var name string
	if start == "client" {
		name = "frpc.ini"
	}
	if start == "server" {
		name = "frps.ini"
	}
	cfg, err := ini.Load(mid.Dir + "/config/frp/" + name)
	if err != nil {
		_, _ = os.Create(mid.Dir + "/config/frp/" + name)
		cfg, _ = ini.Load(mid.Dir + "/config/frp/" + name)
		common, _ := cfg.NewSection("common")
		if mid.GetMode() == "client" {
			_, _ = common.NewKey("server_addr", "x.x.x.x")
			_, _ = common.NewKey("server_port", "7000")
		}
		if mid.GetMode() == "server" {
			_, _ = common.NewKey("bind_port", "7000")
		}
		err = cfg.SaveTo(mid.Dir + "/config/frp/" + name)
		if err != nil {
			mid.Log.Error(fmt.Sprintf("err:%v", err))
		}
	}
	config, err := os.ReadFile(mid.Dir + "/config/frp/" + name)
	if err != nil {
		mid.Log.Error(fmt.Sprintf("err:%v", err))
	}
	return config
}

// UpdateFrp 更新FRP版本
func UpdateFrp(start string) bool {
	var (
		frp Frp
	)
	frp = GetFrpNew(start)
	info := mod.FrpVersion()
	if info == "" {
		info = "Info: Frp New Install"
	}
	if frp.BrowserDownloadUrl == "" {
		mid.Log.Info(mid.RunFuncName() + ":Frp接口请求异常")
		return false
	}
	if "v"+info != frp.TagName {
		url := fmt.Sprintf("%v%v", "https://mirror.ghproxy.com/", frp.BrowserDownloadUrl)
		if con.Down(url, frp.Path) {
			mid.Log.Info("Frp " + frp.TagName + " 下载成功：" + frp.PathName)
		}
		pid := mod.ForPid("frp")
		pids := strings.Split(pid, "\n")
		if pids[0] == "" {
			if !mod.IsExists(mid.Dir + "/config/frp") {
				err := os.MkdirAll(mid.Dir+"/config/frp", 0755)
				if err != nil {
					return false
				}
			}
			err := mod.ExecCommand("tar -zvxf " + mid.Dir + "/config/frp.gz -C " + mid.Dir + "/config/frp")
			if err != nil {
				mid.Log.Error(err.Error())
			}
			if start == "client" {
				err = mod.ExecCommand("cd " + mid.Dir + "/config/frp/frp_* && mv frpc ../")
				if err != nil {
					mid.Log.Error(err.Error())
				}
				err = mod.ExecCommand("rm -rf " + mid.Dir + "/config/frp/frp_*")
				if err != nil {
					mid.Log.Error(err.Error())
				}
			}
			if start == "server" {
				err = mod.ExecCommand("cd " + mid.Dir + "/config/frp/frp_* && mv frps ../")
				if err != nil {
					mid.Log.Error(err.Error())
				}
				err = mod.ExecCommand("rm -rf " + mid.Dir + "/config/frp/frp_*")
				if err != nil {
					mid.Log.Error(err.Error())
				}
			}
		} else {
			if !mod.IsExists(mid.Dir + "/config/frp") {
				err := os.MkdirAll(mid.Dir+"/config/frp", 0755)
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
			err := mod.ExecCommand("tar -zvxf " + mid.Dir + "/config/frp.gz -C " + mid.Dir + "/config/frp")
			if err != nil {
				mid.Log.Error(err.Error())
			}
			if start == "client" {
				err = mod.ExecCommand("cd " + mid.Dir + "/config/frp/frp_* && mv frpc ../")
				if err != nil {
					mid.Log.Error(err.Error())
				}
				err = mod.ExecCommand("rm -rf " + mid.Dir + "/config/frp/frp_*")
				if err != nil {
					mid.Log.Error(err.Error())
				}
			}
			if start == "server" {
				err = mod.ExecCommand("cd " + mid.Dir + "/config/frp/frp_* && mv frps ../")
				if err != nil {
					mid.Log.Error(err.Error())
				}
				err = mod.ExecCommand("rm -rf " + mid.Dir + "/config/frp/frp_*")
				if err != nil {
					mid.Log.Error(err.Error())
				}
			}
		}
		err := mod.ExecCommand("chmod -R 755 " + mid.Dir + "/config/frp")
		if err != nil {
			mid.Log.Error(err.Error())
		}
		return true
	} else {
		return false
	}
}

// GetFrpNew 获取对应操作系统系统的Nps信息
func GetFrpNew(start string) (Frp Frp) {
	var (
		gitApi GitHubReleasesApi
		url    string
	)
	url = "https://api.github.com/repos/fatedier/frp/releases/latest"
	frpInfo := con.GetApi(url, nil)
	err := json.Unmarshal(frpInfo, &gitApi)
	if err != nil {
		mid.Log.Error("FRP线上接口请求受限：" + err.Error())
		return
	}
	architecture, system := mid.GetArchitecture()
	Frp.Path = mid.Dir + "/config/frp.gz"
	Frp.TagName = gitApi.TagName
	Frp.UpdatedAt = mod.Time(gitApi.CreatedAt)
	for _, s := range gitApi.Assets {
		name := strings.Split(s.Name, "_")
		start := strings.Split(name[3], ".")
		if start[0] == architecture && name[2] == system {
			Frp.BrowserDownloadUrl = s.BrowserDownloadUrl
			Frp.PathName = s.Name
		}
	}
	Frp.Starts = start
	Frp.Name = gitApi.Name
	return
}
