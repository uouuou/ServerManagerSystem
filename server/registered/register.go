package registered

import (
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/util"
)

var db = util.GetDB()

// Register 注册程序，客户端主动连接并注册该服务器
type Register struct {
	mod.Model
	Userid     string `json:"userid" gorm:"unique"`                //用户传入时自带的userid 具有唯一性 定义为：当前时间戳+ServerManager+4为随机数字的md5值，且设置为唯一的可用ID如果用户被删除必须换ID才可以注册
	ClientIp   string `json:"client_ip"`                           //接入客户端的外网IP
	Version    string `json:"version"`                             //接入客户端程序版本号
	NatAuth    int    `json:"nat_auth" gorm:"default:2"`           //设置允许连接到本机的NPS和FRP内网穿透上：1为允许连接 2为不允许，默认不允许连接(2)
	NpsConfig  string `json:"nps_config" gorm:"default:8024:2333"` //NPS的服务器端口和密钥，默认端口8024密钥2333采用英文：隔开
	FrpConfig  string `json:"frp_config" gorm:"default:null"`      //frp的配置json
	CronAuth   int    `json:"cron_auth" gorm:"default:2"`          //1为允许连接2为不允许，默认不允许连接
	FrpVersion string `json:"frp_version"`                         //客戶端FRP版本
	NpsVersion string `json:"nps_version"`                         //客戶端NPS版本
	Remark     string `json:"remark"`                              //备注
	Status     bool   `json:"status" gorm:"-"`                     //客户端在线状态
}

//Register 客户端注册
func (r Register) Register(c *gin.Context) {
	var register Register
	err := c.BindJSON(&r)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if r.ClientIp == "" || r.Userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if db.Model(&r).Where("userid = ? and deleted_at IS NULL", r.Userid).Find(&register).RowsAffected <= 0 {
		db.Create(&r)
		mid.DataOk(c, r, "新增成功")
	} else {
		db.Where("userid = ? and deleted_at IS NULL", r.Userid).Updates(&r)
		mid.DataOk(c, r, "更新完成")
	}
}

//RegisterRpc 客户端注册(通过RPC实现)
func RegisterRpc(r Register) mid.ResultBody {
	var register Register
	if r.ClientIp == "" || r.Userid == "" {
		return mid.RpcDataNot(r, r.Userid+":数据错误")
	}
	if db.Model(&r).Where("userid = ? and deleted_at IS NULL", r.Userid).Find(&register).RowsAffected <= 0 {
		if err := db.Create(&r).Error; err == nil {
			return mid.RpcDataOk(r, r.Userid+":写入正常")
		} else {
			return mid.RpcDataErr(err, r.Userid+":写入异常")
		}
	} else {
		db.Where("userid = ? and deleted_at IS NULL", r.Userid).Updates(&r)
		return mid.RpcDataOkUp(register, r.Userid+":更新正常")
	}
}

// List 查询客户端列表
func (r Register) List(c *gin.Context) {
	var register []Register
	pages, Db := mid.GetPages(db, c.Query("page"), c.Query("page_size"), &register)
	if err := Db.Where("deleted_at IS NULL").Find(&register).Error; err == nil {
		for _, s := range mid.GetIDList() {
			for i, registers := range register {
				if registers.FrpConfig == "" {
					register[i].FrpConfig = "null"
				}
				if s == registers.Userid {
					register[i].Status = true
				}
			}
		}
		mid.DataPageOk(c, pages, register, "查询成功")
	} else {
		mid.DataErr(c, err, "查询失败")
	}
}

// MyList 根据userid查询客户端信息
func (r Register) MyList(c *gin.Context) {
	var register Register
	myId := c.Query("userid")
	if myId == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if err := db.Where("userid = ? deleted_at IS NULL", myId).Find(&register).Error; err == nil {
		if register.FrpConfig == "" {
			register.FrpConfig = "null"
		}
		mid.DataOk(c, register, "查询成功")
	} else {
		mid.DataErr(c, err, "查询失败")
	}
}

// Set 设置客户端可用功能
func (r Register) Set(c *gin.Context) {
	var register Register
	err := c.BindJSON(&r)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if r.ClientIp == "" || r.Userid == "" || r.NatAuth == 0 || r.CronAuth == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if r.NatAuth == 1 {
		if r.FrpConfig == "" && r.NpsConfig == "" {
			mid.ClientBreak(c, nil, "未设置NAT客户端参数")
			return
		}
	}
	if db.Model(&r).Where("userid = ? and deleted_at IS NULL", r.Userid).Find(&register).RowsAffected <= 0 {
		mid.DataNot(c, nil, "没有该客户端")
		return
	} else {
		if err = db.Where("userid = ? and deleted_at IS NULL", r.Userid).Updates(&r).Error; err != nil {
			mid.DataErr(c, err, "设置失败")
		} else {
			mid.DataOk(c, r, "设置完成")
		}
	}
}

// Del 删除客户端
func (r Register) Del(c *gin.Context) {
	var register Register
	err := c.BindJSON(&r)
	if err != nil {
		mid.ClientErr(c, err, "格式错误")
		return
	}
	if r.ID == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if db.Model(&r).Where("id = ? and deleted_at IS NULL", r.ID).Find(&register).RowsAffected <= 0 {
		mid.DataNot(c, nil, "没有该客户端")
		return
	} else {
		if err = db.Where("id = ? and deleted_at IS NULL", r.ID).Delete(&r).Error; err != nil {
			mid.DataErr(c, err, "删除失败")
		} else {
			mid.DataOk(c, gin.H{
				"id": r.ID,
			}, "删除成功")
		}
	}
}

// NpsConf 读取客户端NPS设置
func (r Register) NpsConf(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if db.Model(&r).Where("userid = ? and deleted_at IS NULL", userid).Find(&r).RowsAffected <= 0 {
		mid.DataNot(c, nil, "没有该客户端")
		return
	} else {
		mid.DataOk(c, r.NpsConfig, "NPS配置查询成功")
	}
}

// FrpConf 读取客户端FRP设置
func (r Register) FrpConf(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	if db.Model(&r).Where("userid = ? and deleted_at IS NULL", userid).Find(&r).RowsAffected <= 0 {
		mid.DataNot(c, nil, "没有该客户端")
		return
	} else {
		mid.DataOk(c, r.FrpConfig, "FRP配置查询成功")
	}
}
