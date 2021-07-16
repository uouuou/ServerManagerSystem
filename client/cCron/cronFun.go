package cCron

import (
	"github.com/robfig/cron/v3"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/crons"
)

// CronRunList 正在运行的Cron列表
type CronRunList struct {
	EntryID cron.EntryID
	MyCron  crons.MyCron
}

var CronRunLists []CronRunList
var RunCron = cron.New(cron.WithSeconds())

// AddFun 新增Cron的方法
func (m *CronRunList) AddFun(c crons.MyCron) (entryId cron.EntryID, res string, err error) {
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
func (m *CronRunList) RemoveFun(c CronRunList) bool {
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
func (m *CronRunList) Start() bool {
	RunCron.Start()
	return true
}

// Stop 停止Cron
func (m *CronRunList) Stop() bool {
	RunCron.Stop()
	return true
}
