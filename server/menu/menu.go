package menu

import (
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/user"
	"github.com/uouuou/ServerManagerSystem/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var db = util.GetDB()

// Menu 菜单结构体
type Menu struct {
	Id         uint   `json:"id"`
	MenuCode   int    `json:"menu_code" gorm:"primarykey"`   //ID
	MenuName   string `json:"menu_name"`                     //menu名称
	ParentCode int    `json:"parent_code"`                   // 父级ID
	Url        string `json:"url"`                           //路径
	Icon       string `gorm:"type:varchar(20);" json:"icon"` // 图标
	Sort       int    `json:"sort"`                          // 排序值
	Authority  int    `json:"authority" gorm:"default:1"`    //权限：1为可用，2为不可用
}

// TreeList 定义一个序列化数据的结构体
type TreeList struct {
	Id         uint   `json:"id"`
	MenuCode   int    `gorm:"primarykey" json:"menu_code"`   //ID
	MenuName   string `json:"menu_name"`                     //menu名称
	ParentCode int    `json:"parent_code"`                   // 父级ID
	Url        string `json:"url"`                           //路径
	Icon       string `gorm:"type:varchar(20);" json:"icon"` // 图标
	Sort       int    `json:"sort"`                          // 排序值
	Authority  int    `json:"authority" gorm:"default:1"`    //权限：1为可用，2为不可用
	Children   []Menu `json:"children"`                      //子节点
	UpdateUser string `json:"update_user"`                   // 更新人
}

// GetMenuList 查看菜单列表
func GetMenuList(c *gin.Context) {
	var (
		role  user.Role
		users mod.User
	)
	if err := db.Model(&users).Where("name = ?", mid.GetTokenName(c)).Find(&users).Error; err != nil {
		mid.DataErr(c, err, "用户查询失败")
		return
	}
	if users.ID == 1 || users.RoleID == 0 {
		users.RoleID = 1
	}
	if err := db.Model(&user.Role{}).Where("id = ?", users.RoleID).Find(&role).Error; err != nil {
		mid.DataErr(c, err, "角色查询失败")
		return
	}
	treeLists := GetRoleMenu(users, role)
	mid.DataOk(c, treeLists, "success")
}

// GetRoleMenu 获取对应的菜单列表
func GetRoleMenu(users mod.User, role user.Role) []TreeList {
	var (
		menus     []Menu
		treeLists []TreeList
	)
	switch users.RoleID {
	case 1:
		db.Where("authority = 1 and parent_code = 0 and deleted_at IS NULL").Order("sort").Group("menu_code").Find(&menus)
		for _, v := range menus {
			var treeList TreeList
			var menu []Menu
			if v.MenuCode != 1 {
				db.Where("menu_code = ? and parent_code != 0 and authority = 1 and deleted_at IS NULL", v.MenuCode).Order("sort").Find(&menus)
				for _, s := range menus {
					menu = append(menu, s)
				}
			}
			treeList.Icon = v.Icon
			treeList.Id = v.Id
			treeList.MenuCode = v.MenuCode
			treeList.MenuName = v.MenuName
			treeList.ParentCode = v.ParentCode
			treeList.Sort = v.Sort
			treeList.Url = v.Url
			treeList.Children = menu
			treeList.Authority = v.Authority
			treeLists = append(treeLists, treeList)
		}
	default:
		roles := strings.Split(role.RoleCode, ":")
		for _, s := range roles {
			db.Where("authority = 1 and deleted_at IS NULL").Order("sort").Group("menu_code").Find(&menus)
			for _, m := range menus {
				ss, _ := strconv.Atoi(s)
				if m.Id == uint(ss) && m.ParentCode == 0 {
					db.Where("authority = 1 and parent_code = 0 and id = ? and deleted_at IS NULL", ss).Order("sort").Group("menu_code").Find(&menus)
					for _, v := range menus {
						var treeList TreeList
						var menu []Menu

						if v.MenuCode != 1 {
							db.Where("menu_code = ? and parent_code != 0 and authority = 1 and deleted_at IS NULL", v.MenuCode).Order("sort").Find(&menus)
							for _, rs := range roles {
								for _, s := range menus {
									ss, _ := strconv.Atoi(rs)
									if s.Id == uint(ss) {
										menu = append(menu, s)
									}
								}
							}
						}
						treeList.Icon = v.Icon
						treeList.Id = v.Id
						treeList.MenuCode = v.MenuCode
						treeList.MenuName = v.MenuName
						treeList.ParentCode = v.ParentCode
						treeList.Sort = v.Sort
						treeList.Url = v.Url
						treeList.Children = menu
						treeList.Authority = v.Authority
						treeLists = append(treeLists, treeList)
					}
				}
			}
		}
	}
	return treeLists
}

// GetMenuLists 查看菜单列表涵盖被禁用
func GetMenuLists(c *gin.Context) {
	var (
		menus     []Menu
		treeLists []TreeList
		pages     mid.Pages
	)
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	Db := db.Model(&menus)
	pages.TotalAmount = Db.Where("parent_code = 0 and deleted_at IS NULL").Find(&menus).RowsAffected
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
	Db.Where("parent_code = 0 and deleted_at IS NULL").Order("sort").Group("menu_code").Find(&menus)
	for _, v := range menus {
		var treeList TreeList
		var menu []Menu
		if v.MenuCode != 1 {
			db.Where("menu_code = ? and parent_code != 0  and deleted_at IS NULL", v.MenuCode).Order("sort").Find(&menus)
			for _, s := range menus {
				menu = append(menu, s)
			}
		}
		treeList.Icon = v.Icon
		treeList.Id = v.Id
		treeList.MenuCode = v.MenuCode
		treeList.MenuName = v.MenuName
		treeList.ParentCode = v.ParentCode
		treeList.Sort = v.Sort
		treeList.Url = v.Url
		treeList.Children = menu
		treeList.Authority = v.Authority
		treeLists = append(treeLists, treeList)
	}
	mid.DataPageOk(c, pages, treeLists, "success")
}

// AddMenu 新增菜单
func AddMenu(c *gin.Context) {
	var (
		menus mod.Menu
	)
	err := c.ShouldBind(&menus)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if menus.MenuName == "" || menus.Url == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	menus.UpdateUser = mid.GetTokenName(c)
	menuFind := db.Model(&menus).Where("menu_name = ? and deleted_at IS NULL", menus.MenuName).Find(&menus)
	if menuFind.RowsAffected == 0 {
		db.Create(&menus)
		mid.DataOk(c, nil, "新增成功")
	} else {
		mid.DataNot(c, nil, "不可重复添加菜单")
	}
}

// EditMenu 修改菜单
func EditMenu(c *gin.Context) {
	var (
		menus mod.Menu
		menu  mod.Menu
	)
	err := c.ShouldBind(&menus)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if menus.MenuName == "" || menus.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	menus.UpdateUser = mid.GetTokenName(c)
	pushMenu := menus
	menuFind := db.Model(&menus).Where("id = ? and deleted_at IS NULL", menus.ID).Find(&menu)
	if menu.ID != menus.ID {
		mid.DataNot(c, nil, "参数越界")
		return
	}
	if menuFind.RowsAffected == 1 {
		db.Model(&menus).Where("id = ? and deleted_at IS NULL", menus.ID).Updates(&pushMenu)
		mid.DataOk(c, gin.H{
			"name": pushMenu.MenuName,
			"id":   pushMenu.ID,
			"url":  pushMenu.Url,
		}, "修改成功")
	} else {
		mid.DataNot(c, nil, "该菜单不存在")
	}

}

// DelMenu 删除菜单
func DelMenu(c *gin.Context) {
	var (
		menus mod.Menu
	)
	menuId := c.Query("id")
	err := c.ShouldBind(&menus)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if menuId == "" && menus.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if menus.MenuCode == 1 {
		mid.DataNot(c, nil, "不可删除默认菜单")
		return
	} else {
		menuId = strconv.Itoa(int(menus.ID))
	}
	menus.UpdateUser = mid.GetTokenName(c)
	MenuFind := db.Model(&menus).Where("id = ? and deleted_at IS NULL", menuId).Find(&menus)
	if MenuFind.RowsAffected == 1 {
		db.Model(&menus).Where("id = ?", menuId).Delete(&menus)
		mid.DataOk(c, gin.H{
			"name": menus.MenuName,
			"id":   menus.ID,
			"url":  menus.Url,
		}, "删除成功")
	} else {
		mid.DataNot(c, nil, "该菜单不存在")
	}
}
