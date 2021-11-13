package upload

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Upload struct {
	mod.Model
	UploadUser string `json:"upload_user"` //上传者
	RawName    string `json:"raw_name"`    //原生名称
	FileName   string `json:"push_name"`   //公开名称
	Url        string `json:"url"`         //文件地址
	ClientIp   string `json:"client_ip"`   //上传者IP
	UpdateUser string `json:"update_user"` //更新人
}

var db = util.GetDB()

// FilesUpload 文件上传
func (u Upload) FilesUpload(c *gin.Context) {
	var (
		filePath string
		fileUrl  string
		fileName string
		up       Upload
	)
	// 单文件
	file, err := c.FormFile("file")
	if err != nil {
		mid.DataErr(c, err, "No file")
		return
	}
	_dir := mid.Dir + "/upload/" + time.Now().Format("20060102")
	exist, err := PathExists(_dir)
	if err != nil {
		mid.Log.Error(fmt.Sprintf("mkdir failed![%v]", err))
		return
	}

	if !exist {
		err := os.MkdirAll(_dir, os.ModePerm)
		if err != nil {
			mid.Log.Error(fmt.Sprintf("mkdir failed![%v]", err))
		}
	}

	// 上传文件至指定目录
	property := strings.Split(file.Filename, ".")
	if len(property) <= 1 {
		filePath = fmt.Sprintf("%v/%v", _dir, mod.Md5V(strconv.FormatInt(time.Now().Unix(), 10)))
		fileUrl = fmt.Sprintf("%v/%v", "/upload/"+time.Now().Format("20060102"), mod.Md5V(strconv.FormatInt(time.Now().Unix(), 10)))
		fileName = fmt.Sprintf("%v", mod.Md5V(strconv.FormatInt(time.Now().Unix(), 10)))
	} else {
		filePath = fmt.Sprintf("%v/%v.%v", _dir, mod.Md5V(strconv.FormatInt(time.Now().Unix(), 10)), property[1])
		fileUrl = fmt.Sprintf("%v/%v.%v", "/upload/"+time.Now().Format("20060102"), mod.Md5V(strconv.FormatInt(time.Now().Unix(), 10)), property[1])
		fileName = fmt.Sprintf("%v.%v", mod.Md5V(strconv.FormatInt(time.Now().Unix(), 10)), property[1])
	}

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		mid.DataErr(c, err, "文件上传失败")
		return
	}
	up.Url = fileUrl
	up.UploadUser = mid.GetTokenName(c)
	up.FileName = fileName
	up.RawName = file.Filename
	up.ClientIp = c.ClientIP()
	up.UpdateUser = mid.GetTokenName(c)
	if err = db.Model(&up).Create(&up).Error; err != nil {
		mid.DataErr(c, err, "数据写入异常")
	} else {
		mid.DataOk(c, up, "上传成功")
	}
}

// List 文件列表
func (u Upload) List(c *gin.Context) {
	var (
		us []Upload
	)
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &us)
	if err := Db.Find(&us).Error; err != nil {
		mid.DataErr(c, err, "数据查询出错")
		return
	}
	mid.DataPageOk(c, pages, us, "查询成功")
}

// Del 文件删除
func (u Upload) Del(c *gin.Context) {
	var up Upload
	err := c.BindJSON(&u)
	if err != nil {
		mid.ClientErr(c, err, "数据绑定错误")
		return
	}
	ups := db.Where("id = ?", u.ID).Find(&up)
	if ups.Error == nil && ups.RowsAffected >= 0 {
		err = os.Remove(mid.Dir + up.Url)
		if err != nil {
			mid.DataErr(c, err, "文件删除错误")
			return
		}
		if err = db.Where("id = ?", u.ID).Delete(&up).Error; err != nil {
			mid.DataErr(c, err, "文件删除错误")
			return
		}
		mid.DataOk(c, up, "文件删除成功")
	}
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
