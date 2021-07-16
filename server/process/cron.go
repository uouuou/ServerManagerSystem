package process

import (
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	c "github.com/uouuou/ServerManagerSystem/server/crons"
	"strconv"
)

// AutoRun 自动运行进程守护列表
func AutoRun() {
	var p []mod.Process
	var pNow mod.Process
	var pf Process
	taskId, err := c.Crons.AddFunc("*/10 * * * * *", func() {
		if mid.Mode == "server" {
			mod.NewRoutine(c.AutoRunCron)
		}
		db.Model(&p).Where("deleted_at IS NULL").Scan(&p)
	PGo:
		for _, pList := range p {
			if len(mid.AppRunStatus) < len(p) {
				var info mid.AppRunStart
				for _, s := range p {
					info.Name = s.Name
					mid.AppRunStatus = append(mid.AppRunStatus, info)
				}
			}
			if pList.AutoRun == 1 {
				pidOld := mod.ForPidString(pList.RunCmd)
				if pidOld == "" {
					for i := 0; i < 4; i++ {
						pidOld = mod.ForPidString(pList.RunCmd)
					}
				}
				if (mid.Mode == "client" && mid.GetNatAuth() == 2) || (mid.Mode == "client" && mid.GetNatAuth() == 0) {
					continue PGo
				}
				for _, s := range mid.AppRunStatus {
					if s.Name == pList.Name {
						s.Status = false
						if s.Num != 0 {
							if s.Num >= pList.Num {
								continue PGo
							}
						}
					}
				}
				if pidOld == "" || pidOld != pList.Pid {
					run, _ := pf.AddRun(pList.RunCmd, pList.RunPath, pList.Num, pList.Name)
					if !run {
						for i, s := range mid.AppRunStatus {
							if s.Name == pList.Name {
								mid.AppRunStatus[i].Status = false
								mid.AppRunStatus[i].Msg = true
								mid.AppRunStatus[i].Num++
							}
						}
						mid.Log().Error(pList.Name + ":没能正常启动")
						continue PGo
					} else {
						pNow.Pid = mod.ForPidString(pList.RunCmd)
						if pNow.Pid == "" {
							for i, s := range mid.AppRunStatus {
								if s.Name == pList.Name {
									mid.AppRunStatus[i].Status = false
									mid.AppRunStatus[i].Msg = true
									mid.AppRunStatus[i].Num++
								}
							}
							mid.Log().Error(pList.Name + ":没能正常启动")
							continue PGo
						} else {
							db.Model(&p).Where("id = ? and deleted_at IS NULL", pList.ID).Updates(&pNow)
							for i, s := range mid.AppRunStatus {
								if s.Name == pList.Name {
									mid.AppRunStatus[i].Status = true
									mid.AppRunStatus[i].Msg = true
									mid.AppRunStatus[i].Num++
								}
							}
							mid.Log().Info("Process" + pList.Name + ":启动成功 PID " + pNow.Pid)
							continue PGo
						}
					}
				}
				if pidOld == pList.Pid {
					for i, s := range mid.AppRunStatus {
						if s.Name == pList.Name {
							if !s.Msg {
								mid.Log().Info("Process" + pList.Name + ":主程序异常退出,守护进程不受影响 PID " + pidOld)
							}
							mid.AppRunStatus[i].Msg = true
							break
						}
						continue
					}
				}
			}
		}
	})
	if err != nil {
		mid.Log().Error(err.Error())
	}
	c.Crons.Start()
	mid.Log().Info("Process on */10 * * * * * TaskID:" + strconv.Itoa(int(taskId)))
}
