package cNat

import (
	"github.com/uouuou/ServerManagerSystem/client/cProcess"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/nat"
	"strings"
	"time"
)

// NpsInstall 客户端NPS安装
func NpsInstall() {
	var ppp []mod.Process
	var pf cProcess.Process
	nps := nat.GetNpsNew(mid.GetMode())
	info := mod.NpsVersion()
	if info == "" {
		info = "Info: Nps New Install"
		if nat.UpdateNps(mid.GetMode()) {
			mid.Log.Info("安装NPS成功")
		} else {
			mid.Log.Info("安装NPS失败，请重启重试......")
		}
		return
	}
	npsInfo := strings.Split(info, ": ")
	if "v"+npsInfo[1] != nps.TagName {
		if nat.UpdateNps(mid.GetMode()) {
			mid.Log.Info("更新NPS成功")
		} else {
			mid.Log.Info("更新NPS失败，请重启重试......")
		}
		return
	}
	mid.NPSConfig = string(nat.MoNpsConfig(mid.GetMode()))
	oldConf := mid.GetNPSConfig()
	time.Sleep(time.Second * 4)
	nowConf := string(nat.MoNpsConfig(mid.GetMode()))
	db.Model(&ppp).Where("deleted_at IS NULL").Find(&ppp)
	for _, p := range ppp {
		if !mid.FileExist(mid.Dir+"/config/nps/nps.conf") && mid.GetNPSConfig() == "null" {
			if p.Name == "NPS" {
				ps := mod.Process{
					Model: mod.Model{
						ID: p.ID,
					},
				}
				pf.DelRpcProcess(ps)
			}
		}
		if oldConf != nowConf {
			if p.Name == "NPC" {
				ps := mod.Process{
					Model: mod.Model{
						ID: p.ID,
					},
				}
				pf.DelRpcProcess(ps)
			}
		}
	}
	if mid.GetMode() == "server" && mid.GetNPSConfig() != "null" {
		p := mod.Process{
			Name:    "NPS",
			RunPath: mid.Dir + "/config/nps",
			RunCmd:  mid.Dir + "/config/nps/nps",
			Num:     4,
			AutoRun: 1,
			Remark:  "NPS系统自动添加",
		}
		res := pf.AddRpcProcess(p)
		if res.Code == 2000 {
			mid.Log.Info(mid.RunFuncName() + res.Message)
		} else {
			mid.Log.Error(mid.RunFuncName() + res.Message)
		}
	}
	if mid.GetMode() == "client" && mid.GetFRTPConfig() != "null" {
		p := mod.Process{
			Name:    "NPC",
			RunPath: mid.Dir + "/config/nps",
			RunCmd:  mid.Dir + "/config/nps/npc -server=" + mid.GetNPSConfig(),
			Num:     4,
			AutoRun: 1,
			Remark:  "NPC系统自动添加",
		}
		res := pf.AddRpcProcess(p)
		if res.Code == 2000 {
			mid.Log.Info(mid.RunFuncName() + res.Message)
		} else {
			mid.Log.Error(mid.RunFuncName() + res.Message)
		}
	}
	mid.Log.Info(mid.RunFuncName() + ":NPS无需安装或更新......")
}
