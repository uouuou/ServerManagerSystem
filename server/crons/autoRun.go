package crons

import (
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"strings"
)

// 外部构建一个oldCronList来做对比
var oldCronList []CronRunList

// AutoRunCron 开机自启动增加Cron
func AutoRunCron() {
	var (
		crons []MyCron
		m     MyCron
	)
	oldCronList = CronRunLists
	db.Where("deleted_at IS NULL").Find(&crons)
	if len(CronRunLists) == 0 {
		for _, myCron := range crons {
			cits := strings.Split(myCron.Effect, ",")
			for _, cit := range cits {
				if cit != "sms" {
					continue
				}
				_, _, err := m.AddFun(myCron)
				if err != nil {
					mid.Log.Errorf("启动脚本异常：%v", err.Error())
				}
			}
		}
	}
	if len(CronRunLists) != len(crons) {
	c:
		for _, myCron := range crons {
			for _, list := range CronRunLists {
				if myCron.CronName == list.MyCron.CronName {
					continue c
				}
			}
			cits := strings.Split(myCron.Effect, ",")
			for _, cit := range cits {
				if cit != mid.GetCUId() {
					continue
				}
				_, _, err := m.AddFun(myCron)
				if err != nil {
					mid.Log.Errorf("启动脚本异常：%v", err.Error())
				}
			}
		}

	}
	if len(CronRunLists) != len(crons) {
	cs:
		for _, list := range CronRunLists {
			for _, myCron := range crons {
				if myCron.CronName == list.MyCron.CronName {
					continue cs
				}
			}
			m.RemoveFun(list)
			mid.Log.Infof("Cron在线更新成功Entries:%v", RunCron.Entries())
		}
	}
	if len(crons) > 0 && len(oldCronList) != len(CronRunLists) {
		m.Start()
		mid.Log.Infof("Cron启动成功Entries:%v", RunCron.Entries())
	}
}
