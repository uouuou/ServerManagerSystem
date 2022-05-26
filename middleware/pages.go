package middleware

import (
	"gorm.io/gorm"
	"strconv"
)

// GetPages 获取当前结构体的pages信息
func GetPages(db *gorm.DB, pageNow string, pageSizeNow string, value interface{}) (Pages, *gorm.DB) {
	var pages Pages
	page, _ := strconv.Atoi(pageNow)
	pageSize, _ := strconv.Atoi(pageSizeNow)
	Db := db.Model(value)
	Db.Count(&pages.TotalAmount)
	if page > 0 && pageSize > 0 {
		Db.Limit(pageSize).Offset((page - 1) * pageSize)
		pages.Page = page
		pages.PageSize = pageSize
	}
	if page > 0 && pageSize == 0 {
		Db.Limit(pageSize).Offset((page - 1) * pageSize)
		pages.Page = page
		pages.PageSize = 15
	}
	if page == 0 && pageSize > 0 {
		Db.Limit(pageSize)
		pages.Page = 1
		pages.PageSize = pageSize
	}
	if page == 0 && pageSize == 0 {
		pages.Page = 1
		pages.PageSize = 15
		Db = Db.Limit(15)
	}
	return pages, Db
}
