package crons

import (
	"github.com/robfig/cron/v3"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"strings"
)

// CronRunList 正在运行的Cron列表
type CronRunList struct {
	EntryID cron.EntryID
	MyCron  MyCron
}

var CronRunLists []CronRunList
var RunCron = cron.New(cron.WithSeconds())

// AddFun 新增Cron的方法
func (m *MyCron) AddFun(c MyCron) (entryId cron.EntryID, res string, err error) {
	var cronRunList CronRunList
	entryId, err = RunCron.AddFunc(c.Cron, func() {
		res, err = mod.ExecCommandWithResult(c.CronUrl)
		if err != nil {
			return
		}
	})
	cronRunList.EntryID = entryId
	cronRunList.MyCron = c
	CronRunLists = append(CronRunLists, cronRunList)
	return
}

// RemoveFun 删除现有Cron的方法
func (m *MyCron) RemoveFun(c CronRunList) bool {
	RunCron.Remove(c.EntryID)
	e := RunCron.Entries()
	for _, entry := range e {
		if entry.ID == c.EntryID {
			return false
		}
	}
	for i, list := range CronRunLists {
		if list.EntryID == c.EntryID {
			CronRunLists = append(CronRunLists[:i], CronRunLists[i+1:]...)
		}
	}
	return true
}

// Start 启动Cron
func (m *MyCron) Start() bool {
	RunCron.Start()
	return true
}

// Stop 停止Cron
func (m *MyCron) Stop() bool {
	RunCron.Stop()
	return true
}

// RpcCronList RPC客户端获取对应的Cron列表
func RpcCronList(cid string) (myCrons []MyCron) {
	var crons []MyCron
	db.Where("deleted_at IS NULL").Find(&crons)
	for _, myCron := range crons {
		cits := strings.Split(myCron.Effect, ",")
		for _, s := range cits {
			if cid != s {
				continue
			}
			myCrons = append(myCrons, myCron)
		}
	}
	return
}
