package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/menu"
	"github.com/uouuou/ServerManagerSystem/server/user"
	"regexp"
	"strconv"
)

// Login 系统登录接口
func Login(c *gin.Context) {
	var (
		users mod.User
		role  user.Role
	)
	err := c.ShouldBindJSON(&users)
	if err != nil {
		mid.DataErr(c, err, "JSON格式错误")
		return
	}
	password := mod.Md5V(users.Password)
	//token, err := c.Cookie("JSESSIONID")
	//if err != nil {
	//	mod.Log().Error(fmt.Sprintf("err:%v", err))
	//}
	//fmt.Println(token)
	//timeNow := time.Now().Format("2006-01-02 15")
	//tokenNow := mod.Md5V(username + "ServerManagerSystem2021" + timeNow)
	loginResult := db.Model(&users).Where("name = ? ", users.Name).Find(&users)
	if loginResult.RowsAffected <= 0 {
		mid.DataNot(c, nil, "用户不存在")
		return
	}
	loginResult = db.Model(&users).Where("name = ? and password = ?", users.Name, password).Find(&users)
	if loginResult.RowsAffected == 0 {
		mid.DataNot(c, nil, "密码错误")
	} else {
		// 生成Token
		uuid := mid.GetUUID()
		tokenString, _ := mid.GenToken(users.Name, users.ID, uuid)
		//tokenName := c.MustGet("token_username").(string)
		//tokenId := c.MustGet("token_id").(string)
		host := c.Request.Host
		matched, _ := regexp.MatchString("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}", host)
		if matched {
			host = ""
		}
		//c.SetCookie("token", tokenString, 3600*2, "/api", host, false, true)
		if users.Avatar == "" {
			users.Avatar = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80"
		}
		err = mid.Set([]byte(uuid), []byte(users.Name), 60*60)
		if err != nil {
			mid.Log().Error(err.Error())
		}
		err = mid.Set([]byte(uuid+"role"), []byte(strconv.Itoa(users.RoleID)), 60*60)
		if err != nil {
			mid.Log().Error(err.Error())
		}

		if err := db.Model(&user.Role{}).Where("id = ?", users.RoleID).Find(&role).Error; err != nil {
			mid.Log().Error(err.Error())
			mid.DataErr(c, err, "角色查询异常")
			return
		}
		treeLists := menu.GetRoleMenu(users, role)
		marshal, err := json.Marshal(treeLists)
		if err != nil {
			mid.Log().Error(err.Error())
			mid.DataErr(c, err, "数据格式化异常")
			return
		}
		err = mid.Set([]byte(uuid+"menu"), marshal, 1024*1024)
		if err != nil {
			mid.Log().Error(err.Error())
			mid.DataErr(c, err, "数据写入缓存异常")
			return
		}
		mid.LoginData(c, users.Name, users.ID, users.Avatar, uuid, tokenString)
	}
}

// Logout 用户登出
func Logout(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		mid.DataErr(c, nil, "数据错误")
		return
	}
	mid.Del([]byte(uuid))
	mid.Del([]byte(uuid + "menu"))
	mid.Del([]byte(uuid + "role"))
	mid.DataOk(c, gin.H{
		"uuid": uuid,
	}, "退出成功")
}
