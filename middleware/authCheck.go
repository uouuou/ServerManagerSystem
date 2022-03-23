package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Menu 菜单结构体
type Menu struct {
	Id         uint   `json:"id"`
	MenuCode   int    `json:"menu_code" gorm:"primarykey"`   //ID
	MenuName   string `json:"menu_name"`                     //menu名称
	ParentCode int    `json:"parent_code"`                   // 父级ID
	Url        string `json:"url"`                           //路径
	Icon       string `gorm:"type:varchar(20);" json:"icon"` // 图标
	Sort       int    `json:"sort"`                          // 排序值
	Authority  int    `json:"authority"`                     //权限：1为可用，2为不可用
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
	Authority  int    `json:"authority"`                     //权限：1为可用，2为不可用
	Children   []Menu `json:"children"`                      //子节点
	UpdateUser string `json:"update_user"`                   // 更新人
}

var AuthCheckMiddleware = authCheck()

// 一个检测路由是否拥有权限的方法，必须按照路由方地址是发布地址的后继才可以
func authCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, _ := c.Get("uuid")
		menu, _ := Get([]byte(uuid.(string) + "menu"))
		role, _ := Get([]byte(uuid.(string) + "role"))
		var treeList []TreeList
		err := json.Unmarshal(menu, &treeList)
		if err != nil {
			Log.Error(err.Error())
			c.Abort()
		}
		url := strings.Split(c.Request.URL.Path, "/")
	a:
		for i, list := range treeList {
			if string(role) == "1" {
				c.Next()
				break a
			}
		b:
			for _, child := range list.Children {
				urls := strings.Split(child.Url, "/")
				if url[len(url)-3] != urls[len(urls)-2] {
					break b
				}
				for _, m := range list.Children {
					urls = strings.Split(m.Url, "/")
					if url[len(url)-2] == urls[len(urls)-1] {
						c.Next()
						break a
					}
				}
			}
			if i == len(treeList)-1 {
				resultBody := ResultBody{
					Code:    4003,
					Data:    nil,
					Message: "无权限",
				}
				c.JSON(http.StatusForbidden, resultBody)
				c.Abort()
				return
			}
		}
	}
}
