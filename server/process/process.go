package process

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var db = util.GetDB()

type Process struct {
}

// AddProcess 添加守护任务
func (Process) AddProcess(c *gin.Context) {
	var (
		p mod.Process
	)
	err := c.ShouldBind(&p)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if p.Name == "" || p.RunPath == "" || p.RunCmd == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if p.Remark == "" {
		mid.ClientErr(c, nil, "说明不能为空")
		return
	}
	p.UpdateUser = mid.GetTokenName(c)
	find := db.Model(&p).Where("name = ? and deleted_at IS NULL", p.Name).Find(&p)
	if find.RowsAffected == 0 {
		db.Create(&p)
		mid.DataOk(c, nil, "新增成功")
	} else {
		mid.DataNot(c, nil, "该程序已经存在")
	}
}

// ProcessList 查看保持列表
func (Process) ProcessList(c *gin.Context) {
	var (
		processLists  []mod.Process
		processListsS []mod.Process
		process       []mod.Process
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &process)
	if err := Db.Where("remark != ? and deleted_at IS NULL", "").Find(&processLists).Error; err != nil {
		mid.Log().Error(err.Error())
	}
	for _, list := range processLists {
		var pid []string
		pid = mod.ForPids(list.RunCmd)
		if len(pid) <= 1 {
			list.Pid = "未启动"
			list.Running = 2
		} else {
			list.Pid = mod.ForPidString(list.RunCmd)
			list.Running = 1
		}
		list.PLog = list.RunPath + "/log/" + list.Name + "_process.log"
		processListsS = append(processListsS, list)
	}
	mid.DataPageOk(c, pages, processListsS, "success")
}

// EditProcess 修改守护进程
func (Process) EditProcess(c *gin.Context) {
	var (
		p  mod.Process
		ps mod.Process
	)
	err := c.ShouldBind(&p)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if p.Name == "" || p.RunPath == "" || p.RunCmd == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	pushP := p
	pushP.UpdateUser = mid.GetTokenName(c)
	shellFind := db.Model(&p).Where("name = ? and deleted_at IS NULL", p.Name).Find(&ps)
	if shellFind.RowsAffected == 1 {
		db.Model(&p).Where("name = ? and deleted_at IS NULL and id = ?", p.Name, p.ID).Updates(&pushP)
		go func() {
			for {
				pid := mod.ForPids(ps.RunCmd)
				if len(pid) >= 1 {
					for _, pp := range pid {
						_, _ = mod.RunCommand("kill -9 " + pp)
						time.Sleep(time.Millisecond * 100)
						continue
					}
				}
				break
			}
		}()
		mid.DataOk(c, gin.H{
			"name": p.Name,
			"path": p.RunPath,
			"cmd":  p.RunCmd,
			"auto": p.AutoRun,
		}, "修改成功")
	} else {
		mid.DataNot(c, nil, "该保持不存在")
	}

}

// DelProcess 删除保持
func (Process) DelProcess(c *gin.Context) {
	var (
		p mod.Process
	)
	err := c.ShouldBind(&p)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if p.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	p.UpdateUser = mid.GetTokenName(c)
	shellFind := db.Model(&p).Where("id = ? and deleted_at IS NULL", p.ID).Find(&p)
	if shellFind.RowsAffected == 1 {
		if err := db.Model(&p).Where("id = ?", p.ID).Delete(&p).Error; err != nil {
			mid.DataErr(c, err, "删除异常")
			return
		} else {
			go func() {
				for {
					pid := mod.ForPids(p.RunCmd)
					if len(pid) >= 1 {
						for _, pp := range pid {
							_, _ = mod.RunCommand("kill -9 " + pp)
							time.Sleep(time.Millisecond * 100)
							continue
						}
					}
					break
				}
			}()
			mid.DataOk(c, gin.H{
				"name": p.Name,
				"path": p.RunPath,
				"cmd":  p.RunCmd,
				"auto": p.AutoRun,
			}, "删除成功")
			return
		}
	} else {
		mid.DataNot(c, nil, "该服务器不存在")
	}
}

// RunProcess 启动守护
func (p Process) RunProcess(c *gin.Context) {
	var (
		ps   mod.Process
		pNow mod.Process
	)
	err := c.ShouldBind(&ps)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if ps.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	shellFind := db.Model(&pNow).Where("id = ? and deleted_at IS NULL", ps.ID).Find(&ps)
	if shellFind.RowsAffected == 1 {
		pid := mod.ForPid(ps.RunCmd)
		if pid != "" {
			mid.DataNot(c, gin.H{
				"name":    ps.Name,
				"path":    ps.RunPath,
				"cmd":     ps.RunCmd,
				"pid":     pid,
				"auto":    ps.AutoRun,
				"running": 2,
			}, "改程序已经启动")
			return
		}
		run, _ := p.AddRun(ps.RunCmd, ps.RunPath, ps.Num, ps.Name)
		pid = mod.ForPid(ps.RunCmd)
		if run != true {
			mid.DataNot(c, gin.H{
				"name":    ps.Name,
				"path":    ps.RunPath,
				"cmd":     ps.RunCmd,
				"auto":    ps.AutoRun,
				"pid":     pid,
				"running": 1,
			}, "程序启动异常")
		} else {
			pNow.Pid = pid
			db.Model(&pNow).Where("id = ? and deleted_at IS NULL", ps.ID).Updates(&pNow)
			mid.DataOk(c, gin.H{
				"name":    ps.Name,
				"path":    ps.RunPath,
				"cmd":     ps.RunCmd,
				"auto":    ps.AutoRun,
				"pid":     pid,
				"running": 2,
			}, "执行成功")
		}
	} else {
		mid.DataNot(c, nil, "该守护不存在")
	}
}

// OffProcess 关闭
func (p Process) OffProcess(c *gin.Context) {
	var (
		ps   mod.Process
		pNow mod.Process
	)
	err := c.ShouldBind(&ps)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if ps.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	DelAppRunProcess(ps)
	shellFind := db.Model(&pNow).Where("id = ? and deleted_at IS NULL", ps.ID).Find(&ps)
	if shellFind.RowsAffected == 1 {
		pid := mod.ForPid(ps.RunCmd)
		if pid == "" {
			mid.DataNot(c, gin.H{
				"name":    ps.Name,
				"path":    ps.RunPath,
				"cmd":     ps.RunCmd,
				"auto":    ps.AutoRun,
				"running": 1,
			}, "进程未启动")
			return
		} else {
			for {
				pids := mod.ForPids(ps.RunCmd)
				if len(pids) >= 1 {
					for _, pp := range pids {
						_, _ = mod.RunCommand("kill -9 " + pp)
						time.Sleep(time.Millisecond * 100)
						continue
					}
				}
				break
			}
			mid.DataOk(c, gin.H{
				"name":    ps.Name,
				"path":    ps.RunPath,
				"cmd":     ps.RunCmd,
				"auto":    ps.AutoRun,
				"running": 1,
			}, "成功关闭")
		}
	} else {
		mid.DataNot(c, nil, "该守护不存在")
	}
}

// AddRun 运行命令并实时查看运行结果
func (Process) AddRun(command string, path string, num int, name string) (b bool, pid int) {
	err := mod.ExecCommand(fmt.Sprintf("cd %v", path))
	if err != nil {
		return false, pid
	}
	_ = os.MkdirAll(path+"/log", 0755)
	ch := make(chan int, num)
	go func() {
		var nums int
		for i := 0; i < num; i++ {
			pid, err = ExecCommand(command, path, name)
			if err != nil && num != nums {
				nums++
				ch <- nums
			}
		}
	}()
	return true, pid
}

// ExecCommand 运行命令并实时查看运行结果
func ExecCommand(command string, path string, name string) (pid int, err error) {
	cmd := exec.Command("bash", "-c", command)
	//设置写入pgid为pid
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		mid.Log().Error(fmt.Sprintf("%v Error:The command is err: %v", mid.RunFuncName(), err))
		return cmd.Process.Pid, err
	}

	ch := make(chan string, 100)
	stdoutScan := bufio.NewScanner(stdout)
	stderrScan := bufio.NewScanner(stderr)
	go func() {
		for stdoutScan.Scan() {
			line := stdoutScan.Text()
			ch <- line
		}
	}()
	go func() {
		for stderrScan.Scan() {
			line := stderrScan.Text()
			ch <- line
		}
	}()
	go func() {
		err = cmd.Wait()
		close(ch)
	}()
	for line := range ch {
		s := line + "\n"
		logfile := path + "/log/" + name + "_process.log"
		f, err := os.OpenFile(logfile, os.O_WRONLY, 0644)
		if err != nil {
			// 打开文件失败处理
			f, _ = os.Create(logfile)
		}
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, 2)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(s), n)
		f.Close()

	}
	return cmd.Process.Pid, err
}

// DelAppRunProcess 清除被移除队列的运行状态
func DelAppRunProcess(p mod.Process) bool {
	for i, runStatus := range mid.AppRunStatus {
		if runStatus.Name == p.Name {
			mid.AppRunStatus = append(mid.AppRunStatus[:i], mid.AppRunStatus[i+1:]...)
		}
	}
	return true
}
