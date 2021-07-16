package cCron

import (
	"github.com/hprose/hprose-golang/v3/io"
	"github.com/hprose/hprose-golang/v3/rpc"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/server/crons"
	"time"
)

type CronClient struct {
	RpcCronList func(cid string) ([]crons.MyCron, error) `name:"RpcCronList"`
}

// AutoCronRun 启动Cron
func AutoCronRun(client *rpc.Client) {
	var stub *CronClient
	var m CronRunList
	var oldCronList []CronRunList
	io.RegisterName("myCron", (*crons.MyCron)(nil))
	client.UseService(&stub)
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		oldCronList = CronRunLists
		nowCrons, err := stub.RpcCronList(mid.GetCUId())
		if err != nil {
			mid.Log().Error(mid.RunFuncName() + ": RPC调度异常：" + err.Error())
			continue
		}
		if len(CronRunLists) == 0 {
			for _, myCron := range nowCrons {
				_, _, err := m.AddFun(myCron)
				if err != nil {
					mid.Log().Errorf("启动脚本异常：%v", err.Error())
				}
			}
		}
		if len(CronRunLists) != len(nowCrons) {
		c:
			for _, myCron := range nowCrons {
				for _, list := range CronRunLists {
					if myCron.CronName == list.MyCron.CronName {
						continue c
					}
				}
				_, _, err := m.AddFun(myCron)
				if err != nil {
					mid.Log().Errorf("启动脚本异常：%v", err.Error())
				}
			}

		}
		if len(CronRunLists) != len(nowCrons) {
		cs:
			for _, list := range CronRunLists {
				for _, myCron := range nowCrons {
					if myCron.CronName == list.MyCron.CronName {
						continue cs
					}
				}
				m.RemoveFun(list)
				mid.Log().Infof("Cron在线更新成功Entries:%v", RunCron.Entries())
			}
		}
		if len(nowCrons) > 0 && len(oldCronList) != len(nowCrons) {
			m.Start()
			mid.Log().Infof("Cron启动成功Entries:%v", RunCron.Entries())
			continue
		}
	}
}
