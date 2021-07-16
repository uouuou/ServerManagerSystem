package registered

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
)

type SqlRegister struct {
	mod.Model
	Name       string `json:"name"`                             //  连接SQL的命名
	DbType     string `json:"dbType"`                           //  数据库连接 mysql mssql
	DbName     string `json:"dbName"`                           //  数据库名称
	DbUser     string `json:"dbUser"`                           //  数据库用户
	DbPass     string `json:"dbPass"`                           //  数据库密码
	DbHost     string `json:"dbHost"`                           //  对应的数据库地址（被调度的机器所能通信的）
	DbPort     int    `json:"dbPort"`                           //  数据库端口
	Remark     string `json:"remark"`                           //  备注
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"` //更新人
}

// SqlList 查询所有的SQL资产列表
func (s SqlRegister) SqlList(c *gin.Context) {
	var sr []SqlRegister
	var srs []SqlRegister
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &srs)
	if err := Db.Where("deleted_at IS NULL").Find(&sr).Error; err == nil {
		for i := range sr {
			sr[i].DbPass = ""
		}
		mid.DataPageOk(c, pages, sr, "查询成功")
	} else {
		mid.DataErr(c, err, "查询失败")
	}
}

// MySqlList 查询单独某一个资产信息
func (s SqlRegister) MySqlList(c *gin.Context) {
	var sr SqlRegister
	id := c.Query("id")
	if err := db.Where("id = ？ and deleted_at IS NULL", id).Find(&sr).Error; err == nil {
		sr.DbPass = ""
		mid.DataOk(c, sr, "查询成功")
	} else {
		mid.DataErr(c, err, "查询失败")
	}
}

// AddSql 新增sql资产
func (s SqlRegister) AddSql(c *gin.Context) {
	var sqlRegister SqlRegister
	err := c.BindJSON(&s)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if s.DbHost == "" || s.Name == "" || s.DbPort == 0 || s.DbPass == "" || s.DbType == "" || s.DbUser == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if s.DbType != "mysql" && s.DbType != "mssql" {
		mid.DataNot(c, nil, "SQL类型不支持")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	if db.Model(&s).Where("name = ? and deleted_at IS NULL", s.Name).Find(&sqlRegister).RowsAffected > 0 {
		mid.DataNot(c, nil, "该数据库已存在")
		return
	} else {
		if err = db.Create(&s).Error; err != nil {
			mid.DataErr(c, err, "新增失败")
		} else {
			mid.DataOk(c, s, "新增成功")
		}
	}
}

// EditSql 修改SQL资产
func (s SqlRegister) EditSql(c *gin.Context) {
	var sqlRegister SqlRegister
	err := c.BindJSON(&s)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if s.DbHost == "" || s.Name == "" || s.DbPort == 0 || s.DbPass == "" || s.DbType == "" || s.DbUser == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if s.DbType != "sqllite" && s.DbType != "mysql" && s.DbType != "mssql" {
		mid.DataNot(c, nil, "SQL类型不支持")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	if db.Model(&s).Where("id = ? and deleted_at IS NULL", s.ID).Find(&sqlRegister).RowsAffected <= 0 {
		mid.DataNot(c, nil, "不存在该数据")
		return
	} else {
		if err = db.Where("id = ? and deleted_at IS NULL", s.ID).Updates(&s).Error; err != nil {
			mid.DataErr(c, err, "修改失败")
		} else {
			mid.DataOk(c, s, "修改成功")
		}
	}
}

// DelSql 删除SQL资产
func (s SqlRegister) DelSql(c *gin.Context) {
	var sqlRegister SqlRegister
	err := c.BindJSON(&s)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if s.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	s.UpdateUser = mid.GetTokenName(c)
	if db.Model(&s).Where("id = ? and deleted_at IS NULL", s.ID).Find(&sqlRegister).RowsAffected <= 0 {
		mid.DataNot(c, nil, "没有该客户端")
		return
	} else {
		if err = db.Where("id = ? and deleted_at IS NULL", s.ID).Delete(&s).Error; err != nil {
			mid.DataErr(c, err, "删除失败")
		} else {
			mid.DataOk(c, gin.H{
				"id": s.ID,
			}, "删除成功")
		}
	}
}
