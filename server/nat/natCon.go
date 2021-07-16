package nat

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
)

var db = util.GetDB()

// GetNatConf 获取两个工具的可用性评估
func GetNatConf(c *gin.Context) {
	var cn mod.ConfigNat
	fun := c.Query("fun")
	if fun == "" {
		mid.DataNot(c, nil, "参数错误")
		return
	}
	db.Model(&cn).Find(&cn)
	if fun == "frp" {
		mid.DataOk(c, gin.H{
			"id":                 cn.ID,
			"frp_run":            cn.FrpRun,
			"frp_version":        "v" + mod.FrpVersion(),
			"frp_online_version": GetFrpNew(mid.GetMode()).TagName,
		}, "查询成功")
	} else if fun == "nps" {
		mid.DataOk(c, gin.H{
			"id":                 cn.ID,
			"nps_run":            cn.NpsRun,
			"nps_version":        mod.NpsVersion(),
			"nps_online_version": GetNpsNew(mid.GetMode()).TagName,
		}, "查询成功")
	} else {
		mid.DataNot(c, nil, "没有该方法")
	}

}

// SetNatConf 设置NAT功能
func SetNatConf(c *gin.Context) {
	var cn mod.ConfigNat
	err := c.ShouldBind(&cn)
	fun := c.Query("fun")
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if fun == "" {
		mid.DataNot(c, nil, "参数错误")
		return
	}
	cn.UpdateUser = mid.GetTokenName(c)
	if fun == "frp" {
		if err = db.Where("id = ? and deleted_at IS NULL", cn.ID).Updates(&cn).Error; err == nil {
			mid.DataOk(c, gin.H{
				"frp_run": cn.FrpRun,
			}, "正常修改")
		} else {
			mid.DataErr(c, err, "修改遇到异常")
		}
	} else if fun == "nps" {
		if err = db.Where("id = ? and deleted_at IS NULL", cn.ID).Updates(&cn).Error; err == nil {
			mid.DataOk(c, gin.H{
				"nps_run": cn.NpsRun,
			}, "正常修改")
		} else {
			mid.DataErr(c, err, "修改遇到异常")
		}
	} else {
		mid.DataNot(c, nil, "没有该方法")
	}

}
