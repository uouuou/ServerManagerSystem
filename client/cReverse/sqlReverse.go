package cReverse

import (
	"database/sql"
	"encoding/json"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/util"
)

// SqlReverse 客户端反向调用SQL
func SqlReverse(rawSql string) map[string]interface{} {
	var results []map[string]interface{}
	db := util.DB
	sqlRows, err := db.Debug().Raw(rawSql).Rows()
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
