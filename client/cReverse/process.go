package cReverse

import (
	"github.com/uouuou/ServerManagerSystem/client/cProcess"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
)

type CProcess struct {
	Rpc cProcess.Process
}

// Add 客户端远程调度新增守护
func (receiver CProcess) Add(r mod.Process) mid.ResultBody {
	return receiver.Rpc.AddRpcProcess(r)
}

// List 客户端远程调度查看客户端守护列表
func (receiver CProcess) List(p int, ps int) (pLists []mod.Process, pages mid.Pages) {
	return receiver.Rpc.ProcessRpcLists(p, ps)
}

// Del 客户端远程调度删除守护
func (receiver CProcess) Del(r mod.Process) mid.ResultBody {
	return receiver.Rpc.DelRpcProcess(r)
}

// Edit 客户端远程调度修改守护
func (receiver CProcess) Edit(r mod.Process) mid.ResultBody {
	return receiver.Rpc.EditFunProcess(r)
}

// Run 客户端远程调度启动守护
func (receiver CProcess) Run(r mod.Process) mid.ResultBody {
	return receiver.Rpc.RunRpcProcess(r)
}

// Off 客户端远程调度关闭守护
func (receiver CProcess) Off(r mod.Process) mid.ResultBody {
	return receiver.Rpc.OffRpcProcess(r)
}
