package client

import (
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
)

var db = util.GetDB()

// IntSqlStart 同步数据结构，设置基础数据
func IntSqlStart() {
	err := db.AutoMigrate(&mod.Process{})
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + "SQL注册异常" + err.Error())
		return
	}
	err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&mod.Process{})
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + "SQL注册异常" + err.Error())
		return
	}
}
