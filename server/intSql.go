package server

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/crons"
	"github.com/uouuou/ServerManagerSystem/server/net/firewall"
	"github.com/uouuou/ServerManagerSystem/server/net/webshell"
	"github.com/uouuou/ServerManagerSystem/server/registered"
	"github.com/uouuou/ServerManagerSystem/server/upload"
	"github.com/uouuou/ServerManagerSystem/server/user"
	"github.com/uouuou/ServerManagerSystem/util"
)

var db = util.GetDB()

// IntSqlStart 同步数据结构，设置基础数据
func IntSqlStart() {
	var users mod.User
	var menu mod.Menu
	var conf mod.ConfigNat
	var roles user.Role
	_ = db.AutoMigrate(&mod.User{}, &mod.Menu{}, &webshell.ServerInfo{}, &firewall.Firewall{}, &mod.Sub{}, &mod.SubSet{}, &registered.Register{},
		&mod.ConfigNat{}, &mod.Process{}, &mod.Update{}, &registered.SqlRegister{}, &user.Role{}, &upload.Upload{}, &crons.MyCron{})
	_ = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&mod.User{}, &mod.Menu{}, &webshell.ServerInfo{}, &firewall.Firewall{},
		&mod.Sub{}, &mod.SubSet{}, &registered.Register{}, &mod.ConfigNat{}, &mod.Process{}, &mod.Update{}, &registered.SqlRegister{}, &user.Role{},
		&upload.Upload{}, &crons.MyCron{})
	if db.Find(&users).RowsAffected <= 0 {
		users.Name = "admin"
		users.Password = mod.Md5V("123456")
		users.RoleID = 1
		users.Email = "test@email.com"
		db.Create(&users)
	}
	if db.Find(&roles).RowsAffected <= 0 {
		roles.RoleCode = "ADMIN"
		roles.RoleRule = "ADMIN:ALL"
		roles.RoleName = "超级管理员"
		roles.UpdateUser = "系统"
		if err := db.Create(&roles).Error; err != nil {
			mid.Log().Error(err.Error())
		}
	}
	if db.Find(&conf).RowsAffected <= 0 {
		conf.Remarks = "系统默认设置，请勿修改或删除"
		db.Create(&conf)
	}
	if db.Find(&menu).RowsAffected <= 0 {
		menuAdd := fmt.Sprintf("INSERT INTO menus ( menu_code, menu_name, parent_code, url, icon, sort,authority,sys,update_user) VALUES ")
		db.Exec(menuAdd + " (1, '首页', 0, '/', 'el-icon-house', 1, 1, 1, '系统');")
		db.Exec(menuAdd + " (2, '内网穿透', 0, '/nat', 'el-icon-cold-drink', 2, 1, 2, '系统');")
		db.Exec(menuAdd + " (2, 'FRP', 2, '/nat/frp', 'el-icon-goblet', 1, 1, 2, '系统');")
		db.Exec(menuAdd + " (2, 'NPS', 2, '/nat/nps', 'el-icon-goblet-full', 2, 1, 2, '系统');")
		db.Exec(menuAdd + " (3, 'SM图床', 0, '/img', 'el-icon-picture-outline', 3, 1, 1, '系统');")
		db.Exec(menuAdd + " (3, '图床管理', 3, '/img/manage', 'el-icon-camera', 1, 1, 1, '系统');")
		db.Exec(menuAdd + " (3, '图片分享', 5, '/img/share', 'el-icon-share', 2, 1, 1, '系统');")
		db.Exec(menuAdd + " (4, '进程保持', 0, '/process', 'el-icon-takeaway-box', 4, 1, 1, '系统');")
		db.Exec(menuAdd + " (4, '进程保持', 1, '/process/manage', 'el-icon-monitor', 1, 1, 1, '系统');")
		db.Exec(menuAdd + " (5, '网络安全', 0, '/net', 'el-icon-first-aid-kit', 5, 1, 1, '系统');")
		db.Exec(menuAdd + " (5, '访问控制', 5, '/net/access', 'el-icon-coordinate', 1, 1, 1, '系统');")
		db.Exec(menuAdd + " (5, '防火墙', 5, '/net/firewall', 'el-icon-table-lamp', 2, 1, 1, '系统');")
		db.Exec(menuAdd + " (5, 'WebShell', 5, '/net/shell', 'el-icon-connection', 3, 1, 1, '系统');")
		db.Exec(menuAdd + " (5, '定时任务', 5, '/net/cron', 'el-icon-time', 4, 1, 1, '系统');")
		//db.Exec(menuAdd + " (6, '代理设置', 0, '/clash', 'el-icon-connection', 6, 2, 2, '系统');")
		//db.Exec(menuAdd + " (6, '订阅服务', 6, '/clash/sub', 'el-icon-notebook-1', 1, 1, 2, '系统');")
		//db.Exec(menuAdd + " (6, '网易解锁', 6, '/clash/netease', 'el-icon-headset', 2, 1, 2, '系统');")
		//db.Exec(menuAdd + " (6, '节点管理', 6, '/clash/node', 'el-icon-brush', 3, 1, 2, '系统');")
		//db.Exec(menuAdd + " (6, '规则管理', 6, '/clash/rules', 'el-icon-toilet-paper', 4, 1, 2, '系统');")
		//db.Exec(menuAdd + " (6, '代理日志', 6, '/clash/log', 'el-icon-pear', 5, 1, 2, '系统');")
		db.Exec(menuAdd + " (6, '接入管理', 6, '/action/register', 'el-icon-guide', 1, 1, 1, '系统');")
		db.Exec(menuAdd + " (6, '行为管理', 0, '/action', 'el-icon-ice-cream-square', 6, 1, 1, '系统');")
		db.Exec(menuAdd + " (6, 'SQL资产', 6, '/action/sql', 'el-icon-coin', 2, 1, 1, '系统');")
		db.Exec(menuAdd + " (6, '更新管理', 6, '/action/update', 'el-icon-sell', 3, 1, 2, '系统');")
		db.Exec(menuAdd + " (6, '组件编译', 6, '/action/build', 'el-icon-ice-cream-round', 4, 1, 1, '系统');")
		db.Exec(menuAdd + " (7, '系统设置', 0, '/setting', 'el-icon-setting', 7, 1, 2, '系统');")
		db.Exec(menuAdd + " (7, '角色管理', 7, '/setting/role', 'el-icon-office-building', 1, 1, 2, '系统');")
		db.Exec(menuAdd + " (7, '用户设置', 7, '/setting/users', 'el-icon-user', 2, 1, 2, '系统');")
		db.Exec(menuAdd + " (7, '文件管理', 7, '/setting/file', 'el-icon-files', 3, 1, 2, '系统');")
		db.Exec(menuAdd + " (7, '菜单管理', 7, '/setting/menu', 'el-icon-orange', 4, 1, 2, '系统');")
		db.Exec(menuAdd + " (7, '网络设置', 7, '/setting/net', 'el-icon-table-lamp', 5, 1, 2, '系统');")
		db.Exec(menuAdd + " (7, '系统设置', 7, '/setting/system', 'el-icon-shopping-bag-1', 6, 1, 2, '系统');")
		mid.Log().Info("初始化数据库完毕......")
	}
}
