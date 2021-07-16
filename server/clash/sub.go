package clash

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	con "github.com/uouuou/ServerManagerSystem/server"
	"github.com/uouuou/ServerManagerSystem/util"
	"net/url"
)

var db = util.GetDB()

// Sub 订阅数据结构体
type Sub struct {
}

// SubSet 订阅设置
type SubSet struct {
	Path       string `json:"path"`        //订阅文件下载地址
	SubUrl     string `json:"sub_url"`     //订阅服务器
	Include    string `json:"include"`     //指仅保留匹配到的节点，支持正则匹配，需要经过 URLEncode 处理，会覆盖配置文件里的设置
	Exclude    string `json:"exclude"`     //指排除匹配到的节点，支持正则匹配，需要经过 URLEncode 处理，会覆盖配置文件里的设置
	Scv        int    `json:"scv"`         //用于关闭 TLS 节点的证书检查，默认为 false 1：true 2:false
	Udp        int    `json:"udp"`         //用于开启该订阅链接的 UDP，默认为 false 1：true 2:false
	Tls13      int    `json:"tls_13"`      //用于设置是否为节点增加tls1.3开启参数 1：true 2:false
	UpdateUser string `json:"update_user"` // 更新人
}

// GetSubList 订阅列表
func (Sub) GetSubList(c *gin.Context) {
	var (
		sub []mod.Sub
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &sub)
	if err := Db.Where("sub_name != ? and deleted_at IS NULL", "").Find(&sub).Error; err != nil {
		mid.Log().Error(err.Error())
	}
	mid.DataPageOk(c, pages, sub, "success")
}

// AddSub 新增订阅地址
func (Sub) AddSub(c *gin.Context) {
	var (
		s mod.Sub
	)
	err := c.ShouldBind(&s)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if s.SubName == "" && s.SubUrl == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	find := db.Model(&s).Where("sub_name = ? and deleted_at IS NULL", s.SubName).Find(&s)
	if find.RowsAffected == 0 {
		db.Create(&s)
		mid.DataOk(c, s, "新增成功")
	} else {
		mid.DataNot(c, s, "该端订阅已添加")
	}
}

// EditSub 修改订阅地址
func (Sub) EditSub(c *gin.Context) {
	var (
		s mod.Sub
	)
	err := c.ShouldBind(&s)
	if err != nil {
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	if s.SubName == "" && s.SubUrl == "" {
		mid.ClientErr(c, nil, "参数错误")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	find := db.Model(&s).Where("id = ? and deleted_at IS NULL", s.ID).Find(&s)
	if find.RowsAffected != 0 {
		db.Model(&s).Where("sub_name = ? and id = ? and deleted_at IS NULL", s.SubName, s.ID).Updates(&s)
		mid.DataOk(c, s, "修改成功")
	} else {
		mid.DataNot(c, s, "不存在该条订阅")
	}
}

// DelSub 删除订阅地址
func (Sub) DelSub(c *gin.Context) {
	var (
		s mod.Sub
	)
	err := c.ShouldBind(&s)
	if err != nil {
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	if s.SubName == "" && s.SubUrl == "" {
		mid.ClientErr(c, nil, "参数错误")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	find := db.Model(&s).Where("id = ? and deleted_at IS NULL", s.ID).Find(&s)
	if find.RowsAffected != 0 {
		db.Model(&s).Where("id = ?", s.ID).Delete(&s)
		mid.DataOk(c, s, "删除成功")
	} else {
		mid.DataNot(c, s, "不存在该条订阅")
	}
}

// SubNow 立即订阅
func (s SubSet) SubNow(c *gin.Context) {
	var (
		sub    []mod.Sub
		subUrl string
		scv    bool
		udp    bool
		tls    bool
	)
	db.Model(&s).Find(&s)
	db.Model(&sub).Where("deleted_at IS NULL").Find(&sub)
	s.Path = mid.Dir + "/config/clash.yaml"
	for i, l := range sub {
		if i == 0 {
			subUrl = url.QueryEscape(l.SubUrl)
		} else {
			subUrl = subUrl + "|" + url.QueryEscape(l.SubUrl)
		}
	}
	switch s.Scv {
	case 1:
		scv = true
	case 2:
		scv = false
	}
	switch s.Udp {
	case 1:
		udp = true
	case 2:
		udp = false
	}
	switch s.Tls13 {
	case 1:
		tls = true
	case 2:
		tls = false
	}
	subUrlNow := fmt.Sprintf("%vtarget=clashr&url=%v&include=%v&exclude=%v&scv=%v&udp=%v&tls13=%v", s.SubUrl, subUrl, s.Include, s.Exclude, scv, udp, tls)
	if con.Down(subUrlNow, s.Path) {
		ReadConfig()
		go func() {
			Run()
		}()
		mid.DataOk(c, nil, "订阅完成")
	} else {
		mid.DataNot(c, nil, "订阅异常")
	}
}

// SubSet 订阅设置
func (s SubSet) SubSet(c *gin.Context) {
	var subSet SubSet
	err := c.ShouldBind(&subSet)
	if err != nil {
		mid.ClientBreak(c, nil, "格式错误")
		return
	}
	if subSet.SubUrl == "" {
		mid.ClientErr(c, nil, "参数错误")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	find := db.Find(&s)
	if find.RowsAffected == 0 {
		subSetMap := structs.Map(&subSet)
		db.Model(&subSet).Create(subSetMap)
		mid.DataOk(c, subSet, "设置成功")
	} else {
		//通过结构体变量更新字段值, gorm库会忽略零值字段。就是字段值等于0, nil, "", false这些值会被忽略掉，不会更新。如果想更新零值，可以使用map类型替代结构体。
		subSetMap := structs.Map(&subSet)
		db.Model(&subSet).Where("id = 1").Updates(subSetMap)
		mid.DataOk(c, subSet, "修改成功")
	}
}

// GetSubSet 订阅列表
func (s SubSet) GetSubSet(c *gin.Context) {
	db.Model(&s).Find(&s)
	mid.DataOk(c, s, "success")
}
