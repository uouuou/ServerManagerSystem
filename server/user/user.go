package user

import (
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

var db = util.GetDB()

// AddUser 新增用户
func AddUser(c *gin.Context) {
	var (
		users mod.User
		user  mod.User
	)
	err := c.ShouldBind(&users)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if users.Name == "" || users.Password == "" || users.RoleID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	users.Password = mod.Md5V(users.Password)
	r, err := SRole(users.RoleID)
	if err != nil {
		mid.ClientErr(c, err, "角色查询异常")
		return
	}
	users.Role = r.RoleName
	users.UpdateUser = mid.GetTokenName(c)
	userFind := db.Model(&user).Where("name = ? and deleted_at IS NULL", users.Name).Find(&user)
	if userFind.RowsAffected == 0 {
		db.Create(&users)
		mid.DataOk(c, gin.H{
			"name":  users.Name,
			"id":    users.ID,
			"email": users.Email,
			"role":  users.Role,
		}, "新增成功")
	} else {
		mid.DataNot(c, nil, "不可重复添加用户")
	}
}

// ListUser 查看用户列表
func ListUser(c *gin.Context) {
	var (
		users []mod.User
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &users)
	if err := Db.Where("name != ? and deleted_at IS NULL", "").Find(&users).Error; err != nil {
		mid.Log().Error(err.Error())
	}
	for i, user := range users {
		r, err := SRole(user.RoleID)
		if err != nil {
			mid.ClientErr(c, err, "角色查询异常")
			return
		}
		users[i].Role = r.RoleName
		users[i].Password = ""
	}
	mid.DataPageOk(c, pages, users, "查询成功")
}

// EditUser 修改用户
func EditUser(c *gin.Context) {
	var (
		users mod.User
		user  mod.User
	)
	err := c.ShouldBind(&users)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if users.Name == "" || users.ID == 0 || users.RoleID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	pushUser := users
	pushUser.UpdateUser = mid.GetTokenName(c)
	pushUser.Password = mod.Md5V(users.Password)
	r, err := SRole(users.RoleID)
	if err != nil {
		mid.ClientErr(c, err, "角色查询异常")
		return
	}
	pushUser.Role = r.RoleName
	userFind := db.Model(&users).Where("name = ? and deleted_at IS NULL", users.Name).Find(&user)
	if userFind.RowsAffected == 1 {
		db.Model(&users).Where("name = ? and deleted_at IS NULL and id = ?", users.Name, users.ID).Updates(&pushUser)
		mid.DataOk(c, gin.H{
			"name":  pushUser.Name,
			"id":    pushUser.ID,
			"email": pushUser.Email,
			"role":  pushUser.Role,
		}, "修改成功")
	} else {
		mid.DataNot(c, nil, "该用户不存在")
	}

}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	var (
		users mod.User
		user  mod.User
	)
	err := c.ShouldBind(&users)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if users.Name == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	pushUser := users
	pushUser.UpdateUser = mid.GetTokenName(c)
	pushUser.Password = mod.Md5V(users.Password)
	r, err := SRole(users.RoleID)
	if err != nil {
		mid.ClientErr(c, err, "角色查询异常")
		return
	}
	pushUser.Role = r.RoleName
	userFind := db.Model(&users).Where("name = ? and deleted_at IS NULL", users.Name).Find(&user)
	if user.ID != users.ID {
		mid.DataNot(c, nil, "该用户不存在")
		return
	}
	if userFind.RowsAffected == 1 {
		db.Model(&users).Where("name = ? and deleted_at IS NULL", users.Name, users.ID).Updates(&pushUser)
		mid.DataOk(c, gin.H{
			"name":  pushUser.Name,
			"id":    pushUser.ID,
			"email": pushUser.Email,
			"role":  pushUser.Role,
		}, "修改成功")
	} else {
		mid.DataNot(c, nil, "该用户不存在")
	}

}

// DelUser 删除用户
func DelUser(c *gin.Context) {
	var (
		users mod.User
	)
	userId := c.Query("userid")
	err := c.ShouldBind(&users)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if userId == "" && users.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if users.ID == 0 {
		mid.ClientErr(c, nil, "不可删除默认账户")
		return
	} else {
		userId = strconv.Itoa(int(users.ID))
	}
	users.UpdateUser = mid.GetTokenName(c)
	userFind := db.Model(&users).Where("id = ? and deleted_at IS NULL", userId).Find(&users)
	if userFind.RowsAffected == 1 {
		db.Model(&users).Where("id = ?", userId).Delete(&users)
		mid.DataOk(c, gin.H{
			"name":  users.Name,
			"id":    users.ID,
			"email": users.Email,
			"role":  users.Role,
		}, "删除成功")
	} else {
		mid.DataNot(c, nil, "该用户不存在")
	}
}
