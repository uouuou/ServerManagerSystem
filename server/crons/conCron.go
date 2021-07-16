package crons

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
	"strings"
)

var Crons = cron.New(cron.WithSeconds())
var db = util.GetDB()

// MyCron 内部存储Cron的结构体
type MyCron struct {
	mod.Model
	CronName   string `json:"cron_name"`                        //定时任务名称
	Cron       string `json:"cron"`                             //定时任务时间
	CronUrl    string `json:"cron_url"`                         //定时任务脚本地址
	Effect     string `json:"effect" gorm:"type:varchar(1024)"` //作用UUID(UUID需要新增一组默认值为sms的如果为sms则为服务端执行)
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

// MyCronRequest 内部存储Cron的结构体
type MyCronRequest struct {
	mod.Model
	CronName   string   `json:"cron_name"`                        //定时任务名称
	Cron       string   `json:"cron"`                             //定时任务时间
	CronUrl    string   `json:"cron_url"`                         //定时任务脚本地址
	Effect     []string `json:"effect" gorm:"type:varchar(1024)"` //作用UUID(UUID需要新增一组默认值为sms的如果为sms则为服务端执行)
	Effects    string   `json:"effects" gorm:"-"`
	UpdateUser string   `json:"update_user" gorm:"comment:'更新人'"`
}

// List Cron列表
func (m *MyCron) List(c *gin.Context) {
	var (
		crons []MyCron
		req   []MyCronRequest
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &crons)
	if err := Db.Find(&crons).Error; err != nil {
		mid.DataErr(c, err, "数据查询异常")
	}
	for _, cronNow := range crons {
		reqNow := StoSs(c, cronNow)
		req = append(req, reqNow)
	}
	mid.DataPageOk(c, pages, req, "查询成功")
}

//Add 新增Cron
func (m *MyCron) Add(c *gin.Context) {
	var (
		crons   MyCron
		cronNow MyCron
		req     MyCronRequest
	)
	err := c.ShouldBind(&req)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定异常")
		return
	}
	if req.CronName == "" || req.Cron == "" || req.CronUrl == "" {
		mid.DataNot(c, crons, "数据不全请检查核心部分")
		return
	}
	cronFind := db.Model(&crons).Where("cron_name = ? and deleted_at IS NULL", req.CronName).Find(&cronNow)
	crons = SsToS(c, req)
	if cronFind.RowsAffected == 0 {
		if err := db.Create(&crons).Error; err != nil {
			mid.DataErr(c, err, "数据写入异常")
		} else {
			mid.DataOk(c, crons, "新增成功")
		}
	} else {
		mid.DataNot(c, crons, "不可重复添加相同名称的Cron")
	}
}

// Edit 修改Cron
func (m *MyCron) Edit(c *gin.Context) {
	var (
		crons   MyCron
		cronNow MyCron
		req     MyCronRequest
	)
	err := c.ShouldBind(&req)
	if err != nil {
		mid.ClientBreak(c, err, "数据绑定异常")
		return
	}
	if req.CronName == "" || req.Cron == "" || req.CronUrl == "" {
		mid.ClientErr(c, nil, "数据不全请检查核心部分")
		return
	}
	cronFind := db.Model(&crons).Where("cron_name = ? and deleted_at IS NULL", req.CronName).Find(&cronNow)
	crons = SsToS(c, req)
	if cronFind.RowsAffected == 1 {
		if err = db.Model(&crons).Where("cron_name = ?", crons.CronName).Updates(&crons).Error; err != nil {
			mid.DataErr(c, err, "数据写入异常")
			return
		} else {
			mid.DataOk(c, crons, "修改成功")
			return
		}
	} else {
		mid.DataNot(c, crons, "该Cron不存在")
	}
}

// Del 删除Cron
func (m *MyCron) Del(c *gin.Context) {
	var (
		crons   MyCron
		cronNow MyCron
	)
	err := c.ShouldBind(&crons)
	if err != nil {
		mid.ClientBreak(c, err, "数据绑定异常")
		return
	}
	if crons.ID == 0 {
		mid.ClientErr(c, nil, "数据不全请检查核心部分")
		return
	}
	cronFind := db.Model(&crons).Where("id = ? and deleted_at IS NULL", crons.ID).Find(&cronNow)
	crons.UpdateUser = mid.GetTokenName(c)
	if cronFind.RowsAffected == 1 {
		if err := db.Model(&crons).Where("id = ?", crons.ID).Delete(&crons).Error; err != nil {
			mid.DataErr(c, err, "数据删除异常")
			return
		} else {
			mid.DataOk(c, cronNow, "删除成功")
			return
		}
	} else {
		mid.DataNot(c, cronNow, "不存在该Cron")
	}
}

// SsToS 房间数组转换为逗号的字符串
func SsToS(c *gin.Context, req MyCronRequest) MyCron {
	var eff string
	for i, s := range req.Effect {
		if i == 0 {
			eff = s
			continue
		}
		eff = eff + "," + s
	}
	nowCron := MyCron{
		Model: mod.Model{
			ID:        req.ID,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
			DeletedAt: req.DeletedAt,
		},
		CronName:   req.CronName,
		Cron:       req.Cron,
		CronUrl:    req.CronUrl,
		Effect:     eff,
		UpdateUser: mid.GetTokenName(c),
	}
	return nowCron
}

// StoSs 将字符串转换为数组
func StoSs(c *gin.Context, nowCron MyCron) MyCronRequest {
	var req MyCronRequest
	reList := strings.Split(nowCron.Effect, ",")
	for _, s := range reList {
		req.Effect = append(req.Effect, s)
	}
	req.ID = nowCron.ID
	req.CronUrl = nowCron.CronUrl
	req.UpdatedAt = nowCron.UpdatedAt
	req.DeletedAt = nowCron.DeletedAt
	req.CreatedAt = nowCron.CreatedAt
	req.UpdateUser = mid.GetTokenName(c)
	req.Cron = nowCron.Cron
	req.CronName = nowCron.CronName
	req.Effects = nowCron.Effect
	return req
}
