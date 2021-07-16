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
	pages.TotalAmount = Db.Where("deleted_at IS NULL").Find(value).RowsAffected
	if page > 0 && pageSize > 0 {
		Db.Limit(pageSize).Offset((page - 1) * pageSize)
		pages.Page = page
		pages.PageSize = pageSize

	} else if pageSize == -1 {
		pages.Page = page
		pages.PageSize = pageSize
	} else {
		Db = Db.Limit(15)
	}
	return pages, Db
}
