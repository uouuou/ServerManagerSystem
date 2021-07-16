package webshell

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
	"net/http"
	"strconv"
)

var db = util.GetDB()

// ServerInfo 服务器信息
type ServerInfo struct {
	mod.Model
	Memo          string `gorm:"column:memo;size:128;" json:"memo" form:"memo"`                                        // 备注
	ServerAddress string `gorm:"column:server_address;size:128;not null;" json:"server_address" form:"server_address"` // 服务器地址
	UserName      string `gorm:"column:user_name;size:128;not null;" json:"user_name" form:"user_name"`                // 登录用户名
	Password      string `gorm:"column:password;type:char(128);not null;" json:"password" form:"password"`             // 登录密码
	AliasName     string `gorm:"column:alias_name;size:64;" json:"alias_name" form:"alias_name"`                       // 服务器别名
	UpdateUser    string `json:"update_user" gorm:"comment:'更新人'"`                                                     // 更新人
}

// GetShellList 查看服务器列表
func (ServerInfo) GetShellList(c *gin.Context) {
	var (
		serverInfo []ServerInfo
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &serverInfo)
	if err := Db.Where("server_address != ? and deleted_at IS NULL", "").Find(&serverInfo).Error; err != nil {
		mid.Log().Error(err.Error())
	}
	mid.DataPageOk(c, pages, serverInfo, "success")
}

// AddShell 新增服务器
func (ServerInfo) AddShell(c *gin.Context) {
	var (
		serverInfo  ServerInfo
		serverInfos ServerInfo
	)
	err := c.ShouldBind(&serverInfo)
	if err != nil {
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	if serverInfo.ServerAddress == "" || serverInfo.UserName == "" || serverInfo.Password == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	serverInfos.UpdateUser = mid.GetTokenName(c)
	shellFind := db.Model(&serverInfo).Where("server_address = ? and deleted_at IS NULL", serverInfo.ServerAddress).Find(&serverInfos)
	if shellFind.RowsAffected == 0 {
		db.Create(&serverInfo)
		mid.DataOk(c, nil, "新增成功")
	} else {
		mid.DataErr(c, nil, "不可重复添加服务器")
	}
}

// EditShell 修改服务器
func (ServerInfo) EditShell(c *gin.Context) {
	var (
		serverInfo  ServerInfo
		serverInfos ServerInfo
	)
	err := c.ShouldBind(&serverInfo)
	if err != nil {
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	if serverInfo.ServerAddress == "" || serverInfo.UserName == "" || serverInfo.Password == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	pushServerInfo := serverInfo
	pushServerInfo.UpdateUser = mid.GetTokenName(c)
	shellFind := db.Model(&serverInfo).Where("id = ? and deleted_at IS NULL", serverInfo.ID).Find(&serverInfos)
	if shellFind.RowsAffected == 1 {
		db.Model(&serverInfo).Where("deleted_at IS NULL and id = ?", serverInfo.ID).Updates(&pushServerInfo)
		mid.DataOk(c, gin.H{
			"server_address": serverInfo.ServerAddress,
			"id":             serverInfo.ID,
			"alias_name":     serverInfo.AliasName,
		}, "修改成功")
	} else {
		mid.DataErr(c, nil, "该服务器不存在")
	}

}

// DelShell 删除服务器
func (ServerInfo) DelShell(c *gin.Context) {
	var (
		serverInfo ServerInfo
	)
	menuId := c.Query("id")
	err := c.ShouldBind(&serverInfo)
	if err != nil {
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	if menuId == "" && serverInfo.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	menuId = strconv.Itoa(int(serverInfo.ID))
	serverInfo.UpdateUser = mid.GetTokenName(c)
	shellFind := db.Model(&serverInfo).Where("id = ? and deleted_at IS NULL", menuId).Find(&serverInfo)
	if shellFind.RowsAffected == 1 {
		db.Model(&serverInfo).Where("id = ?", menuId).Delete(&serverInfo)
		mid.DataOk(c, gin.H{
			"server_address": serverInfo.ServerAddress,
			"id":             serverInfo.ID,
			"alias_name":     serverInfo.AliasName,
		}, "删除成功")
	} else {
		mid.DataErr(c, nil, "该服务器不存在")
	}
}

// Xterm 启动Xterm连接ssh
func (ServerInfo) Xterm(c *gin.Context) {
	sid, _ := c.GetQuery("sid")
	if getServerInfo(sid, 100, 50).Addr == "" {
		mid.DataErr(c, nil, "没有该设备")
	} else {
		c.HTML(http.StatusOK, "xterm.html", gin.H{
			"sid": sid,
		})
	}
}

// Ws 开启一个ws传输shell数据
func (ServerInfo) Ws(c *gin.Context) {
	WebSocketHandler(c.Writer, c.Request, checkUserToken, getServerInfo)
}

//检查token的合法性
func checkUserToken(token string) bool {
	if _, err := mid.ParseToken(token); err != nil {
		return false
	} else {
		return true
	}
}

//获取对应id的数据
func getServerInfo(sid string, cols int, rows int) (m SshLoginModel) {
	var (
		serverInfo ServerInfo
	)

	if db.Where("id = ?", sid).Find(&serverInfo).RowsAffected > 0 {
		m.PtyCols = uint32(cols - 50)
		m.PtyRows = uint32(rows - 1)
		m.Addr = serverInfo.ServerAddress
		m.UserName = serverInfo.UserName
		m.Pwd = serverInfo.Password
	}
	return
}
