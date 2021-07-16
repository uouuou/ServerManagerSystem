package clash

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	con "github.com/uouuou/ServerManagerSystem/server"
	"strings"
)

type RunClash struct {
	ExternalUI string `yaml,json:"external-ui,omitempty"`
	ConfigFile string `json:"config_file"`
	WebUrl     string `json:"web_url"`
	Path       string `json:"path"`
}

var path string

func init() {
	path = mod.GetCurrentDirectory()
}

func Runs() {
	Run()
}

// Run 启动或重新启动clash
func Run() bool {
	run := RunClash{
		ExternalUI: fmt.Sprintf("%v/config/dashboard/clash-dashboard-gh-pages", path),
		ConfigFile: fmt.Sprintf("%v/config/configClash.yaml", path),
		WebUrl:     "https://github.com/Dreamacro/clash-dashboard/archive/gh-pages.zip",
	}
	pid := mod.ForPid("clash")
	pids := strings.Split(pid, "\n")
	run.installWeb(run)
	if pids[0] == " " {
		if run.Run(run) {
			return true
		} else {
			return false
		}
	} else {
		for _, s := range pids {
			if s != " " {
				_, _ = mod.ExecCommandWithResult(fmt.Sprintf("kill -9 %v", s))
			} else {
				continue
			}
		}
		if run.Run(run) {
			return true
		} else {
			return false
		}
	}
}

func (R *RunClash) Run(clash RunClash) bool {
	runClash := fmt.Sprintf("clash -ext-ui %v  -f %v", clash.ExternalUI, clash.ConfigFile)
	if mod.ExecCommandNoErr(runClash) {
		return true
	} else {
		return false
	}
}

func (R *RunClash) installWeb(clash RunClash) bool {
	clash.Path = "./config/dashboard.zip"
	url := fmt.Sprintf("%v%v", "https://ghproxy.com/", clash.WebUrl)
	if !con.Down(url, clash.Path) {
		mid.Log().Info("Clash " + clash.Path + " 下载失败")
		return false
	}
	_, err := mod.ExecCommandWithResult("unzip -o -d " + mid.Dir + "/config/dashboard " + mid.Dir + "/config/dashboard.zip")
	if err != nil {
		mid.Log().Error("解压缩失败")
		return false
	}
	err = mod.ExecCommand("rm -rf " + mid.Dir + "/config/dashboard.zip")
	if err != nil {
		mid.Log().Error(err.Error())
		return false
	}
	return true
}
