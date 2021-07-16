package cProcess

import (
	"bufio"
	"fmt"
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

// ProcessRpcLists 查看保持列表
func (Process) ProcessRpcLists(page int, pageSize int) (pLists []mod.Process, pages mid.Pages) {
	var (
		processLists []mod.Process
		processNow   []mod.Process
	)
	Db := db.Model(&processNow)
	pages.TotalAmount = Db.Find(&processNow).RowsAffected
	if page > 0 && pageSize > 0 {
		Db.Limit(pageSize).Offset((page - 1) * pageSize)
		pages.Page = page
		pages.PageSize = pageSize

	} else if pageSize == -1 {
		pages.Page = page
		pages.PageSize = pageSize
	} else {
		Db = Db.Limit(15)
	}
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
		pLists = append(pLists, list)
	}
	return
}

// AddRpcProcess 添加守护任务
func (Process) AddRpcProcess(p mod.Process) mid.ResultBody {
	var (
		ps mod.Process
	)
	if p.Name == "" || p.RunPath == "" || p.RunCmd == "" {
		return mid.RpcDataNot(p, "数据异常")
	}
	if p.Remark == "" {
		return mid.RpcDataNot(p, "备注不可为空")
	}
	find := db.Model(&p).Where("name = ? and deleted_at IS NULL", p.Name).Find(&ps)
	if find.RowsAffected == 0 {
		db.Create(&p)
		return mid.RpcDataOk(p, p.Name+":新增成功")
	} else {
		if ps.Pid != p.Pid {
			if err := db.Model(&p).Where("name = ? and deleted_at IS NULL", p.Name).Updates(&p).Error; err == nil {
				return mid.RpcDataOk(p, p.Name+":更新成功")
			} else {
				return mid.RpcDataNot(p, p.Name+":更新失败")
			}
		}
		return mid.RpcDataOkUp(p, p.Name+":暂无更新")
	}
}

// DelRpcProcess 删除保持
func (Process) DelRpcProcess(p mod.Process) mid.ResultBody {
	if p.ID == 0 {
		return mid.RpcDataNot(p, "数据异常")
	}
	if p.ID == 1 || p.ID == 2 {
		return mid.RpcDataNot(p, "不可删除默认值")
	}
	shellFind := db.Model(&p).Where("id = ? and deleted_at IS NULL", p.ID).Find(&p)
	if shellFind.RowsAffected == 1 {
		if err := db.Model(&p).Where("id = ?", p.ID).Delete(&p).Error; err != nil {
			return mid.RpcDataErr(err, p.Name+":删除成功")
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
			return mid.RpcDataOk(p, p.Name+":删除成功")
		}
	} else {
		return mid.RpcDataNot(p, p.Name+":该守护不存在")
	}
}

// EditFunProcess 修改守护进程
func (Process) EditFunProcess(p mod.Process) mid.ResultBody {
	var (
		ps mod.Process
	)
	if p.Name == "" || p.RunPath == "" || p.RunCmd == "" || p.ID == 0 {
		fmt.Println(p)
		return mid.RpcDataNot(nil, "提交数据异常")
	}
	pushP := p
	shellFind := db.Model(&p).Where("name = ? and deleted_at IS NULL", p.Name).Find(&ps)
	if shellFind.RowsAffected == 1 {
		if err := db.Model(&p).Where("name = ? and deleted_at IS NULL and id = ?", p.Name, p.ID).Updates(&pushP).Error; err != nil {
			return mid.RpcDataErr(err, "修改异常")
		} else {
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
			return mid.RpcDataOk(p, "修改成功")
		}
	} else {
		return mid.RpcDataNot(nil, "该保持不存在")
	}

}

// RunRpcProcess 启动守护
func (p Process) RunRpcProcess(ps mod.Process) mid.ResultBody {
	var pNow mod.Process
	if ps.ID == 0 {
		return mid.RpcDataNot(p, "数据异常")
	}
	shellFind := db.Model(&pNow).Where("id = ? and deleted_at IS NULL", ps.ID).Find(&ps)
	if shellFind.RowsAffected == 1 {
		pid := mod.ForPid(ps.RunCmd)
		ps.Pid = pid
		ps.Running = 2
		if pid != "" {
			return mid.RpcDataOk(ps, "改程序已经启动")
		}
		run, _ := p.AddRun(ps.RunCmd, ps.RunPath, ps.Num, ps.Name)
		pid = mod.ForPid(ps.RunCmd)
		ps.Pid = pid
		ps.Running = 1
		if run != true {
			return mid.RpcDataNot(ps, "程序启动异常")
		} else {
			pNow.Pid = pid
			ps.Pid = pid
			ps.Running = 2
			db.Model(&pNow).Where("id = ? and deleted_at IS NULL", ps.ID).Updates(&pNow)
			return mid.RpcDataOk(ps, "执行成功")
		}
	}
	return mid.RpcDataNot(nil, "该守护不存在")
}

// OffRpcProcess 关闭
func (p Process) OffRpcProcess(ps mod.Process) mid.ResultBody {
	var pNow mod.Process
	if ps.ID == 0 {
		return mid.RpcDataNot(p, "数据异常")
	}
	DelAppRunProcess(ps)
	shellFind := db.Model(&pNow).Where("id = ? and deleted_at IS NULL", ps.ID).Find(&ps)
	if shellFind.RowsAffected == 1 {
		pid := mod.ForPid(ps.RunCmd)
		ps.Pid = pid
		ps.Running = 1
		if pid == "" {
			return mid.RpcDataNot(ps, "进程未启动")
		} else {
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
			return mid.RpcDataOk(ps, "成功关闭")
		}
	} else {
		return mid.RpcDataNot(nil, "该守护不存在")
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
