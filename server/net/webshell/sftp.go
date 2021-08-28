package webshell

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	goSftp "github.com/pkg/sftp"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	goSsh "golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type File struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	IsDir   bool     `json:"isDir"`
	Mode    string   `json:"mode"`
	IsLink  bool     `json:"isLink"`
	ModTime mod.Time `json:"modTime"`
	Size    int64    `json:"size"`
}

type SftpLoginModel struct {
	user    string
	passwd  string
	addr    string
	SshType int
}

// SftpCreate 创建一个sftp连接
func SftpCreate(login SftpLoginModel) (*goSftp.Client, error) {
	config := goSsh.ClientConfig{
		User:            login.user,
		Timeout:         2 * time.Second,
		HostKeyCallback: goSsh.InsecureIgnoreHostKey(),
	}
	if login.SshType == 1 {
		config.Auth = []goSsh.AuthMethod{goSsh.Password(login.passwd)}
	} else if login.SshType == 2 {
		config.Auth = []goSsh.AuthMethod{publicKeyAuthFunc(login.passwd)}
	} else {
		return nil, errors.New("密钥不存在！")
	}
	conn, err := goSsh.Dial("tcp", login.addr, &config)
	if err != nil {
		return nil, err
	}

	c, err := goSftp.NewClient(conn)
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return c, nil
}

// SftpCat 获取sftp文件
func SftpCat(c *gin.Context) {
	var (
		loginInfo SftpLoginModel
		isFile    bool
	)
	sid := c.Query("sid")
	filepath := c.Query("path")
	if filepath == "" || sid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	fileExt := path.Ext(filepath)
	fileTypes := strings.Split(mid.ReadFileType(), "|")
	for _, fileType := range fileTypes {
		if fileType == fileExt {
			isFile = true
		}
	}
	if !isFile {
		mid.DataOk(c, gin.H{
			"text": "ReadType不允许读取该类型",
			"name": "异常提醒",
		}, "ReadType不允许读取该类型")
		return
	}
	loginInfo = getSftpServerInfo(sid)
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	fileInfo, err := create.Stat(filepath)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	if fileInfo.IsDir() {
		mid.DataErr(c, nil, filepath+" 是目录不能查查看文件内容")
		return
	}
	f, err := create.Open(filepath)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		mid.ClientErr(c, nil, "文件读取异常")
		return
	}
	mid.DataOk(c, gin.H{
		"text": string(b),
		"name": fileInfo.Name(),
	}, "读取成功")
}

// SftPUpload 上传一个文件
func SftPUpload(c *gin.Context) {
	var (
		loginInfo SftpLoginModel
	)
	file, err := c.FormFile("file")
	sid := c.PostForm("sid")
	remoteDir := c.PostForm("dir")
	if err != nil || remoteDir == "" || sid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	filename := file.Filename
	src, err := file.Open()
	if err != nil {
		mid.ClientErr(c, nil, "文件打开失败")
		return
	}
	remoteFile := path.Join(remoteDir, filename)
	loginInfo = getSftpServerInfo(sid)
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	dstFile, err := create.Create(remoteFile)
	if err != nil {
		mid.ClientErr(c, err, "参数错误")
		return
	}
	defer func(dstFile *goSftp.File) {
		err := dstFile.Close()
		if err != nil {

		}
	}(dstFile)
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err != io.EOF {
				mid.DataErr(c, err, "上传错误")
			} else {
				break
			}
		}
		_, _ = dstFile.Write(buf[:n])
	}
	mid.DataOk(c, gin.H{
		"file": filename,
		"path": remoteFile,
	}, "上传成功")
	return
}

// SftpDownload 通过sftp下载文件
func SftpDownload(c *gin.Context) {
	var (
		loginInfo SftpLoginModel
	)
	sid := c.Query("sid")
	filePath := c.Query("path")
	sourceType := c.Query("type") //file or dir
	token := c.Query("token")
	// 获取带后缀的文件名称
	//filenameWithSuffix := path.Base(filePath)
	if filePath == "" || sid == "" || token == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	_, err := mid.ParseToken(token)
	if err != nil {
		mid.ClientBreak(c, nil, "认证错误")
		return
	}
	loginInfo = getSftpServerInfo(sid)
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		resultBody := mid.ResultBody{
			Code:    4002,
			Data:    err.Error(),
			Message: "SFTP连接异常",
		}
		c.JSON(http.StatusOK, resultBody)
		return
	}
	if sourceType == "file" {
		fi, err := create.Stat(filePath)
		if err != nil {
			mid.ClientErr(c, err, "文件读取异常")
			return
		}
		f, err := create.Open(filePath)
		defer func(f *goSftp.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)
		if err != nil {
			mid.ClientErr(c, err, "文件读取异常")
			return
		}
		nameString:=strings.Split(f.Name(),"/")
		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, nameString[len(nameString)-1]),
		}
		c.DataFromReader(http.StatusOK, fi.Size(), "application/octet-stream", f, extraHeaders)
		return
	} else if sourceType == "dir" {
		buf := new(bytes.Buffer)
		w := zip.NewWriter(buf)
		err := zipAddFiles(w, create, filePath, "/")
		if err != nil {
			mid.ClientErr(c, err, "文件压缩异常")
			return
		}
		// Make sure to check the error on Close.
		err = w.Close()
		if err != nil {
			mid.ClientErr(c, err, "文件关闭异常")
			return
		}
		dName := time.Now().Format("2006_01_02T15_04_05Z07.zip")
		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, dName),
		}
		c.DataFromReader(http.StatusOK, int64(buf.Len()), "application/zip", buf, extraHeaders)
		return
	} else {
		mid.DataErr(c, nil, "下载异常")
		return
	}
}

// SftpLs 获取对应文件列表
func SftpLs(c *gin.Context) {
	var (
		loginInfo SftpLoginModel
	)
	sid := c.Query("sid")
	remoteDir := c.Query("dir")
	if sid == "" || remoteDir == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	loginInfo = getSftpServerInfo(sid)
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	fileInfos, err := create.ReadDir(remoteDir)
	if err != nil {
		mid.ClientErr(c, err, "SFTP文件夹读取异常")
		return
	}
	var files = make([]File, 0)
	if len(fileInfos) != 0 {
		for i := range fileInfos {

			// 忽略因此文件
			if strings.HasPrefix(fileInfos[i].Name(), ".") {
				continue
			}

			file := File{
				Name:    fileInfos[i].Name(),
				Path:    path.Join(remoteDir, fileInfos[i].Name()),
				IsDir:   fileInfos[i].IsDir(),
				Mode:    fileInfos[i].Mode().String(),
				IsLink:  fileInfos[i].Mode()&os.ModeSymlink == os.ModeSymlink,
				ModTime: mod.Time(fileInfos[i].ModTime()),
				Size:    fileInfos[i].Size(),
			}

			files = append(files, file)
		}
		mid.DataOk(c, files, "读取成功")
	} else {
		mid.DataOk(c, files, "没有了😄")
	}

}

// SftpRm sftp的删除接口
func SftpRm(c *gin.Context) {
	var loginInfo SftpLoginModel
	type rmInfo struct {
		Sid int    `json:"sid"`
		Key string `json:"key"`
	}
	var rmInfos rmInfo
	err := c.ShouldBind(&rmInfos)
	if err != nil {
		mid.Log().Error(fmt.Sprintf("err:%v", err))
		mid.ClientBreak(c, err, "格式错误")
		return
	}
	if rmInfos.Sid == 0 || rmInfos.Key == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	loginInfo = getSftpServerInfo(strconv.Itoa(rmInfos.Sid))
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	stat, err := create.Stat(rmInfos.Key)
	if err != nil {
		mid.DataErr(c, err, "不存在该文件或文件夹")
		return
	}
	if stat.IsDir() {
		fileInfos, err := create.ReadDir(rmInfos.Key)
		if err != nil {
			mid.DataErr(c, err, "文件夹读取异常")
			return
		}

		for i := range fileInfos {
			if err := create.Remove(path.Join(rmInfos.Key, fileInfos[i].Name())); err != nil {
				mid.DataErr(c, err, "文件删除异常")
				return
			}
		}

		if err := create.RemoveDirectory(rmInfos.Key); err != nil {
			mid.DataErr(c, err, "文件夹删除异常")
			return
		}
	} else {
		if err := create.Remove(rmInfos.Key); err != nil {
			mid.DataErr(c, err, "文件夹删除异常")
			return
		}
	}
	mid.DataOk(c, gin.H{
		"rm_info": rmInfos.Key,
	}, "删除成功")
}

// SftpMkdir sftp的创建文件夹接口
func SftpMkdir(c *gin.Context) {
	var loginInfo SftpLoginModel
	sid := c.Query("sid")
	remoteDir := c.Query("dir")
	if sid == "" || remoteDir == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	loginInfo = getSftpServerInfo(sid)
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	if err := create.Mkdir(remoteDir); err != nil {
		mid.DataErr(c, err, "文件夹创建异常")
		return
	}
	mid.DataOk(c, gin.H{
		"mkdir": remoteDir,
	}, "创建成功")
}

// SftpRenameEndpoint sftp的文件或文件夹重命名接口
func SftpRenameEndpoint(c *gin.Context) {
	var loginInfo SftpLoginModel
	sid := c.Query("sid")
	oldName := c.Query("oldName")
	newName := c.Query("newName")
	if sid == "" || oldName == "" || newName == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	loginInfo = getSftpServerInfo(sid)
	create, err := SftpCreate(loginInfo)
	defer func(create *goSftp.Client) {
		err := create.Close()
		if err != nil {
			return
		}
	}(create)
	if err != nil {
		mid.ClientErr(c, err, "SFTP连接异常")
		return
	}
	if err := create.Rename(oldName, newName); err != nil {
		mid.DataErr(c, err, "重命名异常")
		return
	}
	mid.DataOk(c, gin.H{
		"oldName": oldName,
		"newName": newName,
	}, "重命名成功")
}

//获取对应id的数据
func getSftpServerInfo(sid string) (m SftpLoginModel) {
	var (
		serverInfo ServerInfo
	)
	if db.Where("id = ?", sid).Find(&serverInfo).RowsAffected > 0 {
		m.addr = serverInfo.ServerAddress
		m.user = serverInfo.UserName
		m.passwd = serverInfo.Password
		m.SshType = serverInfo.SshType
	}
	return
}

//压缩文件夹使其可以下载
func zipAddFiles(w *zip.Writer, sftpC *goSftp.Client, basePath, baseInZip string) error {
	// Open the Directory
	files, err := sftpC.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("sftp 读取目录 %s 失败:%s", basePath, err)
	}

	for _, file := range files {
		thisFilePath := basePath + "/" + file.Name()
		if file.IsDir() {

			err := zipAddFiles(w, sftpC, thisFilePath, baseInZip+file.Name()+"/")
			if err != nil {
				return fmt.Errorf("递归目录%s 失败:%s", thisFilePath, err)
			}
		} else {

			dat, err := sftpC.Open(thisFilePath)
			if err != nil {
				return fmt.Errorf("sftp 读取文件失败 %s:%s", thisFilePath, err)
			}
			// Add some files to the archive.
			zipElePath := baseInZip + file.Name()
			f, err := w.Create(zipElePath)
			if err != nil {
				return fmt.Errorf("写入zip writer header失败 %s:%s", zipElePath, err)
			}
			b, err := ioutil.ReadAll(dat)
			if err != nil {
				return fmt.Errorf("ioutil read all failed ：%v", err)
			}
			_, err = f.Write(b)
			if err != nil {
				return fmt.Errorf("写入zip writer 内容 bytes失败:%s", err)
			}
		}
	}
	return nil
}
