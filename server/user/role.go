package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"strconv"
	"strings"
)

// Role 角色主结构体
type Role struct {
	mod.Model
	RoleName   string `json:"role_name"`
	RoleCode   string `json:"role_code"`
	RoleRule   string `json:"role_rule"`
	UpdateUser string `json:"update_user"`
}

// RoleRequest 角色数据返回机构体
type RoleRequest struct {
	mod.Model
	RoleName   string `json:"role_name"`
	RoleCode   []int  `json:"role_code"`
	RoleRule   string `json:"role_rule"`
	UpdateUser string `json:"update_user"`
}

// Add 新增角色
func (r Role) Add(c *gin.Context) {
	var req RoleRequest
	err := c.BindJSON(&req)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定错误")
		return
	}
	if req.RoleName == "" || len(req.RoleCode) == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	res := IntToString(c, req)
	if db.Model(&r).Where("role_name == ?", res.RoleName).Find(&r).RowsAffected >= 1 {
		mid.DataNot(c, res, "角色名称不可重复")
		return
	}
	if err = db.Model(&r).Create(&res).Error; err != nil {
		mid.DataErr(c, err, "新增角色错误")
		return
	}
	mid.DataOk(c, res, "新增角色成功")
}

// List 角色列表
func (r Role) List(c *gin.Context) {
	var (
		rs []Role
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &rs)
	if err := Db.Find(&rs).Error; err != nil {
		mid.DataErr(c, err, "查询角色错误")
		return
	}
	res := StringToInt(c, rs)
	mid.DataPageOk(c, pages, res, "查询成功")
}

func (r Role) Edit(c *gin.Context) {
	var req RoleRequest
	err := c.BindJSON(&req)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定错误")
		return
	}
	if req.RoleName == "" || len(req.RoleCode) == 0 || req.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if req.ID == 1 {
		mid.ClientBreak(c, nil, "不可修改超级管理员")
		return
	}
	res := IntToString(c, req)
	if db.Model(&r).Where("role_name == ?", res.RoleName).Find(&r).RowsAffected < 1 {
		mid.DataNot(c, res, "不存在该角色")
		return
	}
	if err = db.Model(&r).Where("id = ?", r.ID).Updates(&res).Error; err != nil {
		mid.DataErr(c, err, "角色更新失败")
		return
	}
	mid.DataOk(c, r, "角色更新成功")
}

// Del 删除角色
func (r Role) Del(c *gin.Context) {
	var req RoleRequest
	var user mod.User
	err := c.BindJSON(&req)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定错误")
		return
	}
	if req.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if req.ID == 1 {
		mid.DataNot(c, nil, "不可删除默认权限")
		return
	}
	if u := db.Model(&user).Where("role_id = ?", req.ID).Find(&user).RowsAffected; u > 0 {
		mid.DataNot(c, nil, "该角色下存在用户，不可删除")
		return
	}
	if err = db.Model(&r).Where("id = ?", req.ID).Delete(&r).Error; err != nil {
		mid.DataErr(c, err, "角色删除失败")
		return
	}
	mid.DataOk(c, r, "角色删除成功")
}

// IntToString 输入的int和存储的string的转换
func IntToString(c *gin.Context, req RoleRequest) (res Role) {
	var user string
	if req.ID == 1 {
		user = "ADMIN"
	} else {
		user = "USER"
	}
	for i, r := range req.RoleCode {
		if i == 0 {
			res.RoleCode = strconv.Itoa(r)
			res.RoleRule = user + strconv.Itoa(r)
			continue
		}
		res.RoleCode = res.RoleCode + ":" + strconv.Itoa(r)
		res.RoleRule = res.RoleRule + ":" + strconv.Itoa(r)
	}
	res.RoleName = req.RoleName
	res.UpdateUser = mid.GetTokenName(c)
	return
}

// StringToInt 输出的时候将存储的string转换为int
func StringToInt(c *gin.Context, res []Role) (req []RoleRequest) {
	var reqList RoleRequest
	var menu []mod.Menu
	db.Model(&menu).Find(&menu)
	for _, re := range res {
		var reCodeList []int
		if re.RoleCode != "ADMIN" {
			reList := strings.Split(re.RoleCode, ":")
			for _, s := range reList {
				r, _ := strconv.Atoi(s)
				reCodeList = append(reCodeList, r)
			}
		} else {
			for _, i := range menu {
				reCodeList = append(reCodeList, int(i.ID))
			}
			re.RoleRule = "ADMIN:ALL"
		}
		reqList.RoleCode = reCodeList
		reqList.ID = re.ID
		reqList.RoleName = re.RoleName
		reqList.Model = re.Model
		reqList.UpdateUser = mid.GetTokenName(c)
		reqList.RoleRule = re.RoleRule
		req = append(req, reqList)
	}
	return
}

// SRole 获取对应ID的权限信息
func SRole(id int) (r Role, err error) {
	if id == 0 {
		return r, errors.New("传入ID异常")
	}
	if err = db.Model(&r).Where("id = ?", id).Find(&r).Error; err != nil {
		return r, err
	} else {
		return r, nil
	}
}
