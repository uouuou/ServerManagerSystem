package cNat

import (
	"github.com/uouuou/ServerManagerSystem/client/cProcess"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/nat"
	"github.com/uouuou/ServerManagerSystem/util"
)

var db = util.GetDB()

// FrpInstall 客户端FRP安装
func FrpInstall() {
	var ppp []mod.Process
	var pf cProcess.Process
	nps := nat.GetFrpNew(mid.GetMode())
	info := mod.FrpVersion()
	if info == "" {
		info = "Info: Nps New Install"
		if nat.UpdateFrp(mid.GetMode()) {
			mid.Log.Info("安装FRP成功")
		} else {
			mid.Log.Info("安装NPS失败，请重启重试......")
		}
		return
	}
	if "v"+info != nps.TagName {
		if nat.UpdateFrp(mid.GetMode()) {
			mid.Log.Info("更新FRP成功")
		} else {
			mid.Log.Info("更新FRP失败，请重启重试......")
		}
		return
	}
	mid.FRPConfig = string(nat.FrpConfig(mid.GetMode()))
	db.Model(&ppp).Where("deleted_at IS NULL").Find(&ppp)
	for _, p := range ppp {
		if !mid.FileExist(mid.Dir+"/config/frp/frps.ini") && mid.GetFRTPConfig() == "null" {
			if p.Name == "FRPS" {
				ps := mod.Process{
					Model: mod.Model{
						ID: p.ID,
					},
				}
				pf.DelRpcProcess(ps)
			}
		}
		if !mid.FileExist(mid.Dir+"/config/frp/frpc.ini") && mid.GetFRTPConfig() == "null" {
			if p.Name == "FRPC" {
				ps := mod.Process{
					Model: mod.Model{
						ID: p.ID,
					},
				}
				pf.DelRpcProcess(ps)
			}
		}
	}
	if mid.GetMode() == "server" && mid.GetFRTPConfig() != "null" {
		p := mod.Process{
			Name:    "FRPS",
			RunPath: mid.Dir + "/config/frp",
			RunCmd:  mid.Dir + "/config/frp/frps -c " + mid.Dir + "/config/frp/frps.ini",
			Num:     4,
			AutoRun: 1,
			Remark:  "FRPS系统自动添加",
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
			Name:    "FRPC",
			RunPath: mid.Dir + "/config/frp",
			RunCmd:  mid.Dir + "/config/frp/frpc -c " + mid.Dir + "/config/frp/frpc.ini",
			Num:     4,
			AutoRun: 1,
			Remark:  "FRPS系统自动添加",
		}
		res := pf.AddRpcProcess(p)
		if res.Code == 2000 {
			mid.Log.Info(mid.RunFuncName() + res.Message)
		} else {
			mid.Log.Error(mid.RunFuncName() + res.Message)
		}
	}
	mid.Log.Info(mid.RunFuncName() + ":FRP无需安装或更新......")
}
