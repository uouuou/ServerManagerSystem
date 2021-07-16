package cReverse

import (
	"database/sql"
	"encoding/json"
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

// SqlConfig 数据库配置
type SqlConfig struct {
	DbType string `json:"dbType" yaml:"dbType"` // 数据库连接 mysql sqlLite
	DbName string `json:"dbName" yaml:"dbName"` // 如果是mysql dbName
	DbUser string `yaml:"DbUser" json:"dbUser"` // 如果是mysql dbUser
	DbPass string `yaml:"DbPass" json:"dbPass"` // 如果是mysql dbPass
	DbHost string `yaml:"DbHost" json:"dbHost"` // 如果是mysql dbHost
	DbPort int    `yaml:"DbPort" json:"dbPort"` // 如果是mysql dbPort
}

// RunSql 运行远程调度SQL的结构体
type RunSql struct {
	SqlConfig SqlConfig `json:"sql_config"`
	SqlRaw    string    `json:"sql_raw"`
	Cid       string    `json:"cid"`
}

// RunSqlOnId 运行远程调度SQL的结构体
type RunSqlOnId struct {
	SqlId  int    `json:"sql_id"`
	SqlRaw string `json:"sql_raw"`
	Cid    string `json:"cid"`
}

// ClientSql 客户端连接SQL服务器
func ClientSql(sqlConfig SqlConfig) (DB *gorm.DB) {
	var (
		sqlErr error
	)
	// 设置数据库
	if sqlConfig.DbType == "" {
		return nil
	}
	switch sqlConfig.DbType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", sqlConfig.DbUser, sqlConfig.DbPass, sqlConfig.DbHost, sqlConfig.DbPort, sqlConfig.DbName)
		DB, sqlErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	case "mssql":
		configSql := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", sqlConfig.DbHost, sqlConfig.DbUser, sqlConfig.DbPass, sqlConfig.DbPort, sqlConfig.DbName)
		DB, sqlErr = gorm.Open(sqlserver.Open(configSql), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	case "sqllite":
		dsn := mid.Dir + "/config/server.db"
		if !models.FileExist(dsn) {
			configDir := mid.Dir + "/config"
			_, err := os.Stat(configDir)
			if err != nil {
				if os.IsNotExist(err) {
					err := os.MkdirAll(configDir, os.ModePerm)
					if err != nil {
						mid.Log().Info(fmt.Sprintf("err:%v", err))
					}
					return
				}
				return
			}
			// 创建一个数据库文件
			file, err := os.Create(dsn)
			if err != nil {
				mid.Log().Info(fmt.Sprintf("err:%v", err))
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
		}
		DB, sqlErr = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	}
	if DB != nil {
		if sqlErr != nil {
			mid.Log().Error(fmt.Sprintf("SqlConfig connect error %v\n", sqlErr))
		}
		if DB.Error != nil {
			mid.Log().Error(fmt.Sprintf("database error %v\n", DB.Error))
		}
		sqlDB, err := DB.DB()
		if err != nil {
			mid.Log().Error(fmt.Sprintf("sql err:%v\n", err))
		}
		// 连接池最多同时打开的连接数
		sqlDB.SetMaxOpenConns(200)
		//连接池里最大空闲连接数
		sqlDB.SetMaxIdleConns(50)
		//连接池里面的连接最大存活时长
		sqlDB.SetConnMaxLifetime(time.Hour * 2)
		//连接池里面的连接最大空闲时长
		sqlDB.SetConnMaxIdleTime(time.Hour)
		return
	} else {
		return nil
	}
}

// SpecifySqlReverse 服务端指定调用某一个SQL服务器的SQL
func SpecifySqlReverse(runSql RunSql) map[string]interface{} {
	var results []map[string]interface{}
	db := ClientSql(runSql.SqlConfig)
	if db == nil {
		return mid.ResNotOk(nil, "数据库连接失败")
	}
	sqlDb, err := db.DB()
	if err != nil {
		mid.Log().Error(fmt.Sprintf("本地SQL准备异常:%v\n", err))
	}
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			return
		}
	}(sqlDb)
	sqlRows, err := db.Debug().Raw(runSql.SqlRaw).Rows()
	if err != nil {
		mid.Log().Error(mid.RunFuncName() + ":数据库执行异常 " + err.Error())
		return mid.ResErr(err, "数据库执行异常")
	}
	cols, err := sqlRows.Columns()
	if err != nil {
		mid.Log().Error(mid.RunFuncName() + ":数据读取数量异常 " + err.Error())
		return mid.ResErr(err, "数据读取数量异常")
	}
	defer func(sqlRows *sql.Rows) {
		err := sqlRows.Close()
		if err != nil {
			return
		}
	}(sqlRows)
	if len(cols) <= 0 {
		return mid.ResNotOk(nil, "未能获取数据")
	}
	for sqlRows.Next() {
		var row = make([]interface{}, len(cols))
		var rows = make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			rows[i] = &row[i]
		}

		err := sqlRows.Scan(rows...)
		if err != nil {
			mid.Log().Error(mid.RunFuncName() + ":数据匹配异常 " + err.Error())
			return mid.ResErr(err, "数据匹配异常")
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			rowMap[col] = row[i]
		}
		results = append(results, rowMap)
	}
	// 序化化后方便数据传输（否则可能导致map[interface{}]interface{}导致JSON序化化失败）
	marshal, err := json.Marshal(results)
	if err != nil {
		return nil
	}
	return mid.ResOk(marshal, "数据正常")
}
