package firewall

import (
	"fmt"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
)

var db = util.GetDB()

type Firewall struct {
	mod.Model
	Action     int    `json:"action"`                           //1:开放端口 2：屏蔽IP
	Port       int    `json:"port"`                             //需要开放的端口号
	Ip         string `json:"ip"`                               //需要屏蔽的IP地址
	Remarks    string `json:"remarks"`                          //说明备注
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"` //更新人
}

type FireList struct {
	Id         uint     `json:"id"`
	UpdatedAt  mod.Time `json:"updated_at"` //更新时间
	Action     string   `json:"action"`     //行为：开放端口[80]
	Status     string   `json:"status"`     //端口使用状态1：使用中 2：未使用
	Remarks    string   `json:"remarks"`    //说明备注
	UpdateUser string   `json:"update_user"`
}

// AddFirewall 新增一个开放的端口
func (Firewall) AddFirewall(c *gin.Context) {
	var (
		fire Firewall
	)
	err := c.ShouldBind(&fire)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if fire.Action == 1 && fire.Port == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if fire.Action == 2 && fire.Ip == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if fire.Remarks == "" {
		mid.ClientErr(c, err, "说明不能为空")
		return
	}
	fire.UpdateUser = mid.GetTokenName(c)
	if fire.Action == 1 {
		find := db.Model(&fire).Where("port = ? and deleted_at IS NULL", fire.Port).Find(&fire)
		if find.RowsAffected == 0 {
			db.Create(&fire)
			mid.DataOk(c, nil, "新增成功")
		} else {
			mid.DataNot(c, nil, "该端口已经开放")
		}
	}
	if fire.Action == 2 {
		find := db.Model(&fire).Where("ip = ? and deleted_at IS NULL", fire.Ip).Find(&fire)
		if find.RowsAffected == 0 {
			db.Create(&fire)
			mid.DataOk(c, nil, "新增成功")
		} else {
			mid.DataNot(c, nil, "该IP已经屏蔽")
		}
	}
}

// FirewallList 查看防火墙列表
func (Firewall) FirewallList(c *gin.Context) {
	var (
		fireList  FireList
		fireLists []FireList
		fire      []Firewall
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &fire)
	if err := Db.Where("remarks != ? and deleted_at IS NULL", "").Find(&fire).Error; err != nil {
		mid.Log().Error(err.Error())
	}
	for _, s := range fire {
		if s.Action == 1 {
			fireList.Action = fmt.Sprintf("放行端口：[%v]", s.Port)
			fireList.Remarks = s.Remarks
			fireList.Id = s.ID
			fireList.UpdatedAt = s.UpdatedAt
			if mod.PortIsUse(s.Port) {
				fireList.Status = "正常"
			} else {
				fireList.Status = "未使用"
			}
		}
		if s.Action == 2 {
			fireList.Action = fmt.Sprintf("阻止访问：[%v]", s.Ip)
			fireList.Remarks = s.Remarks
			fireList.Id = s.ID
			fireList.UpdatedAt = s.UpdatedAt
			fireList.Status = "正常"

		}
		fireLists = append(fireLists, fireList)
	}
	mid.DataPageOk(c, pages, fireLists, "success")
}

// DelFirewall 删除防火墙
func (Firewall) DelFirewall(c *gin.Context) {
	var (
		fire Firewall
	)
	err := c.ShouldBind(&fire)
	if err != nil {
		mid.Log().Error(fmt.Sprintf("err:%v", err))
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if fire.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if fire.ID == 0 {
		mid.DataNot(c, nil, "不可删除默认防火墙")
		return
	}
	fire.UpdateUser = mid.GetTokenName(c)
	id := fire.ID
	userFind := db.Model(&fire).Where("id = ? and deleted_at IS NULL", id).Find(&fire)
	if userFind.RowsAffected == 1 {
		db.Model(&fire).Where("id = ?", id).Delete(&fire)
		mid.DataOk(c, gin.H{
			"id":      fire.ID,
			"remarks": fire.Remarks,
			"action":  fire.Action,
			"port":    fire.Port,
			"ip":      fire.Ip,
		}, "删除成功")
	} else {
		mid.DataNot(c, nil, "该防火墙不存在")
	}
}
