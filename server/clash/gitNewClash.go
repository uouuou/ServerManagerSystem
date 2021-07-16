package clash

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	con "github.com/uouuou/ServerManagerSystem/server"
	"strings"
	"time"
)

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

// Clash 在线clash信息
type Clash struct {
	ClashName          string    `json:"clash_name"`
	TagName            string    `json:"tag_name"`
	Name               string    `json:"name"`
	UpdatedAt          time.Time `json:"updated_at"`
	Path               string    `json:"path"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
}

// OnlineClashVersion 获取线上的clash版本信息
func OnlineClashVersion(c *gin.Context) {
	var (
		clash Clash
	)
	version := c.Query("version")
	if version == "" {
		mid.ClientErr(c, nil, "参数错误")
		return
	}
	clash = GetClashNew(version)
	mid.DataOk(c, clash, "查询成功")
}

// UpdateClashNew 更新clash版本
func UpdateClashNew(c *gin.Context) {
	var (
		clash Clash
	)
	version := c.Query("version")
	if version == "" {
		mid.ClientErr(c, nil, "参数错误")
		return
	}
	go func() {
		clash = GetClashNew(version)
	}()
	if UpdateClash(version) {
		mid.DataOk(c, clash, "更新成功")
	} else {
		mid.DataNot(c, clash, "暂无更新")
	}
}

// UpdateClash 更新clash版本
func UpdateClash(version string) bool {
	var (
		clash Clash
	)
	run := RunClash{
		ExternalUI: fmt.Sprintf("%v/config/dashboard/clash-dashboard-gh-pages", path),
		ConfigFile: fmt.Sprintf("%v/config/configClash.yaml", path),
		WebUrl:     "https://github.com/Dreamacro/clash-dashboard/archive/gh-pages.zip",
	}
	clash = GetClashNew(version)
	info := mod.ClashVersion()
	if info == "" {
		info = "clash New Install"
	}
	clashInfo := strings.Split(info, " ")
	if clashInfo[1] != clash.TagName {
		url := fmt.Sprintf("%v%v", "https://ghproxy.com/", clash.BrowserDownloadUrl)
		if con.Down(url, clash.Path) {
			mid.Log().Info("Clash " + clash.TagName + " 下载成功：" + clash.ClashName)
		}
		pid := mod.ForPid("clash")
		pids := strings.Split(pid, "\n")
		if pids[0] == " " {
			err := mod.ExecCommand("gunzip -c " + mid.Dir + "/config/clash.gz > /usr/local/bin/clash")
			if err != nil {
				mid.Log().Error(err.Error())
			}
			go run.Run(run)
		} else {
			for _, s := range pids {
				if s != " " {
					_, _ = mod.ExecCommandWithResult(fmt.Sprintf("kill -9 %v", s))
				} else {
					continue
				}
			}
			err := mod.ExecCommand("gunzip -c " + mid.Dir + "/config/clash.gz > /usr/local/bin/clash")
			if err != nil {
				mid.Log().Error(err.Error())
			}
			go run.Run(run)
		}

		err := mod.ExecCommand("chmod -R 755 /usr/local/bin/clash")
		if err != nil {
			mid.Log().Error(err.Error())
		}
		return true
	} else {
		return false
	}
}

// GetClashNew 获取对应操作系统系统的clash信息
func GetClashNew(version string) (clash Clash) {
	var (
		gitApi   GitHubReleasesApi
		clashUrl string
	)
	if version == "premium" {
		clashUrl = "https://api.github.com/repos/Dreamacro/clash/releases/tags/premium"
	}
	if version == "openSource" {
		clashUrl = "https://api.github.com/repos/Dreamacro/clash/releases/latest"
	}
	clashInfo := con.GetApi(clashUrl, nil)
	err := json.Unmarshal(clashInfo, &gitApi)
	if err != nil {
		mid.Log().Error(err.Error())
	}
	architecture, system := mid.GetArchitecture()
	clash.Path = mid.Dir + "/config/clash.gz"
	if version == "premium" {
		name := strings.Split(gitApi.Name, " ")
		clash.TagName = name[1]
	}
	if version == "openSource" {
		clash.TagName = gitApi.TagName
	}
	clash.UpdatedAt = gitApi.CreatedAt
	for _, s := range gitApi.Assets {
		name := strings.Split(s.Name, "-")
		if name[1] == system && name[2] == architecture {
			clash.BrowserDownloadUrl = s.BrowserDownloadUrl
			clash.ClashName = s.Name
		}
	}
	clash.Name = gitApi.Name
	return
}
