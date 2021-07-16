package process

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"net/http"
	"strconv"
)

type CProcessFun struct {
}

// ReverseProcess 远程调用客户端方法的结构体
type ReverseProcess struct {
	Add  func(r mod.Process) (mid.ResultBody, error)
	List func(p int, ps int) ([]mod.Process, mid.Pages, error)
	Edit func(r mod.Process) (mid.ResultBody, error)
	Del  func(r mod.Process) (mid.ResultBody, error)
	Run  func(r mod.Process) (mid.ResultBody, error)
	Off  func(r mod.Process) (mid.ResultBody, error)
}

// List 获取客户端的进程保持列表
func (CProcessFun) List(c *gin.Context) {
	var proxy *ReverseProcess
	userid := c.Query("userid")
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	mid.Caller.UseService(&proxy, userid)
	list, pages, err := proxy.List(page, pageSize)
	if err != nil {
		mid.ClientErr(c, err, "客户端守护查询失败")
		return
	}
	mid.DataPageOk(c, pages, list, "数据查询成功")
}

// Add 对指定客户端新增守护
func (CProcessFun) Add(c *gin.Context) {
	var (
		p     mod.Process
		proxy *ReverseProcess
	)
	userid := c.Query("userid")
	err := c.ShouldBind(&p)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if p.Name == "" || p.RunPath == "" || p.RunCmd == "" || userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if p.Remark == "" {
		mid.ClientErr(c, nil, "说明不能为空")
		return
	}
	p.UpdateUser = mid.GetTokenName(c)
	mid.Caller.UseService(&proxy, userid)
	res, err := proxy.Add(p)
	if err != nil {
		mid.ClientErr(c, err, "远程调度出错")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Edit 对指定客户端修改守护
func (CProcessFun) Edit(c *gin.Context) {
	var (
		p     mod.Process
		proxy *ReverseProcess
	)
	userid := c.Query("userid")
	err := c.ShouldBind(&p)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if p.Name == "" || p.RunPath == "" || p.RunCmd == "" || userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	p.UpdateUser = mid.GetTokenName(c)
	mid.Caller.UseService(&proxy, userid)
	res, err := proxy.Edit(p)
	if err != nil {
		mid.ClientErr(c, err, "远程调度出错")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Del 对指定客户端删除守护
func (CProcessFun) Del(c *gin.Context) {
	var (
		p     mod.Process
		proxy *ReverseProcess
	)
	userid := c.Query("userid")
	err := c.ShouldBind(&p)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if p.ID == 0 || userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	p.UpdateUser = mid.GetTokenName(c)
	mid.Caller.UseService(&proxy, userid)
	res, err := proxy.Del(p)
	if err != nil {
		mid.ClientErr(c, err, "远程调度出错")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Run 客户端启动守护
func (CProcessFun) Run(c *gin.Context) {
	var (
		ps    mod.Process
		proxy *ReverseProcess
	)
	userid := c.Query("userid")
	err := c.ShouldBind(&ps)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if ps.ID == 0 || userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	mid.Caller.UseService(&proxy, userid)
	res, err := proxy.Run(ps)
	if err != nil {
		mid.ClientErr(c, err, "远程调度出错")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Off 客户端启动守护
func (CProcessFun) Off(c *gin.Context) {
	var (
		ps    mod.Process
		proxy *ReverseProcess
	)
	userid := c.Query("userid")
	err := c.ShouldBind(&ps)
	if err != nil || userid == "" {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if ps.ID == 0 || userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	mid.Caller.UseService(&proxy, userid)
	res, err := proxy.Off(ps)
	if err != nil {
		mid.ClientErr(c, err, "远程调度出错")
		return
	}
	c.JSON(http.StatusOK, res)
}
