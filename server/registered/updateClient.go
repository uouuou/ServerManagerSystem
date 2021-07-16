package registered

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
)

func PushUpdate() (up mod.Update, err error) {
	var updates mod.Update
	if err = db.Model(&updates).Last(&up).Error; err != nil {
		mid.Log().Error("获取更新数据失败")
		return
	}
	return
}

// GetUpdateVersion 获取更新版本信息
func GetUpdateVersion(c *gin.Context) {
	var (
		updates []mod.Update
		updateS []mod.Update
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &updates)
	if err := Db.Find(&updates).Error; err != nil {
		mid.DataErr(c, err, "获取数据出错")
		return
	}
	for i, s := range updates {
		if i+1 == int(pages.TotalAmount) {
			s.Now = true
		}
		updateS = append(updateS, s)
	}
	mid.DataPageOk(c, pages, updateS, "success")
}

// SetUpdateVersion 设置用于更新的版本信息
func SetUpdateVersion(c *gin.Context) {
	var (
		update mod.Update
	)
	err := c.ShouldBind(&update)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定异常")
		return
	}
	update.UpdateUser = mid.GetTokenName(c)
	if err = db.Create(&update).Error; err != nil {
		mid.DataErr(c, err, "数据写入异常")
	} else {
		mid.DataOk(c, update, "新增更新成功")
	}
}

// DelUpdateVersion 删除更新程序
func DelUpdateVersion(c *gin.Context) {
	var (
		update mod.Update
	)
	err := c.ShouldBind(&update)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定异常")
		return
	}
	if update.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	update.UpdateUser = mid.GetTokenName(c)
	userFind := db.Model(&update).Where("id = ? and deleted_at IS NULL", update.ID).Find(&update)
	if userFind.RowsAffected == 1 {
		db.Model(&update).Where("id = ?", update.ID).Delete(&update)
		mid.DataOk(c, gin.H{
			"id": update.ID,
		}, "删除成功")
	} else {
		mid.DataNot(c, nil, "该ID不存在")
	}
}
