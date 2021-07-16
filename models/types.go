package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/hprose/hprose-golang/v3/io"
	"gorm.io/gorm"
	"time"
)

// Time 系统内部的时间格式化
// MarshalJSON规范JSON时涉及到的时间格式
// Scan Gorm的情况下查询时的时间格式化
// Value Gorm情况下写入数据时的时间格式化
// String 当时间是string的情况下格式化时间
// UnmarshalJSON 在 c.ShouldBindJSON 时，会调用 field.UnmarshalJSON 方法
type Time time.Time

func (t Time) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

func (t Time) MarshalJSON() ([]byte, error) {
	if (t == Time{}) {
		formatted := fmt.Sprintf("\"%s\"", "")
		return []byte(formatted), nil
	} else {
		formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
		return []byte(formatted), nil
	}
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = Time(time.Time{})
		return
	}
	// 指定解析的格式
	now, err := time.Parse(`"2006-01-02 15:04:05"`, string(data))
	*t = Time(now)
	return
}

func (t *Time) Scan(v interface{}) error {
	switch vt := v.(type) {
	case string:
		// 字符串转成 time.Time 类型
		tTime, _ := time.Parse("2006-01-02 15:04:05", vt)
		*t = Time(tTime)
	case time.Time:
		*t = Time(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

func init() {
	io.RegisterValueEncoder(Time{}, io.GetValueEncoder(time.Time{}))
	io.RegisterValueDecoder(Time{}, io.GetValueDecoder(time.Time{}))
}

// Config 读取config文件
type Config struct {
	Port string `yaml:"Port"`
	Frpc string `yaml:"Frpc"`
	Num  string `yaml:"Num"`
}

// User 用户基本信息
type User struct {
	Model
	Name       string `json:"name" gorm:"comment:'用户名'"`
	Password   string `json:"password" gorm:"comment:'密码'"`
	Role       string `json:"role" gorm:"comment:'角色名称'"`
	RoleID     int    `json:"role_id" gorm:"comment:'角色id'"`
	Email      string `json:"email" gorm:"comment:'邮件地址'"`
	Avatar     string `json:"avatar" gorm:"comment:'头像地址'"`
	SCKey      string `json:"sc_key" gorm:"comment:'server酱'"`
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

// Menu 系统菜单控制
type Menu struct {
	Model
	MenuCode   int    `json:"menu_code"`                     //主菜单ID
	MenuName   string `json:"menu_name"`                     //menu名称
	ParentCode int    `json:"parent_code"`                   // 父级ID
	Url        string `json:"url"`                           //路径
	Icon       string `gorm:"type:varchar(20);" json:"icon"` // 图标
	Sort       int    `json:"sort"`                          // 排序值
	Sys        int    `json:"sys" gorm:"default:1"`          //系统级 ：1为用户级，2为系统级
	Authority  int    `json:"authority" gorm:"default:1"`    //权限：1为可用，2为不可用
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

// Sub 订阅数据结构体
type Sub struct {
	Model
	SubName    string `json:"sub_name"` //订阅名称
	SubUrl     string `json:"sub_url"`  //订阅地址
	Remarks    string `json:"remarks"`  //说明备注
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

// SubSet 订阅设置
type SubSet struct {
	Model
	Path       string `json:"path"`    //订阅文件下载地址
	SubUrl     string `json:"sub_url"` //订阅服务器
	Include    string `json:"include"` //指仅保留匹配到的节点，支持正则匹配，需要经过 URLEncode 处理，会覆盖配置文件里的设置
	Exclude    string `json:"exclude"` //指排除匹配到的节点，支持正则匹配，需要经过 URLEncode 处理，会覆盖配置文件里的设置
	Scv        int    `json:"scv"`     //用于关闭 TLS 节点的证书检查，默认为 false 1：true 2:false
	Udp        int    `json:"udp"`     //用于开启该订阅链接的 UDP，默认为 false 1：true 2:false
	Tls13      int    `json:"tls_13"`  //用于设置是否为节点增加tls1.3开启参数 1：true 2:false
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

// ConfigNat NAT系统配置表
type ConfigNat struct {
	Model
	NpsRun     int    `json:"nps_run" gorm:"default:2"` //是否启用nps 1为启用 2为不启用
	FrpRun     int    `json:"frp_run" gorm:"default:2"` //是否启用frp 1为启用 2为不启用
	Remarks    string `json:"remarks"`                  //说明备注
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

// Update 客户端更新
type Update struct {
	Model
	Version    string `json:"version" gorm:"comment:'版本号'"`
	UrlLinux   string `json:"urlLinux" gorm:"comment:'arm程序地址'"`
	UrlArm     string `json:"urlArm" gorm:"comment:'arm程序地址'"`
	Remark     string `json:"remark" gorm:"comment:'备注可以不填写'"`
	Now        bool   `json:"now" gorm:"comment:'当前版本'"`
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"`
}

type Process struct {
	Model
	Name       string `json:"name"`                             //守护名称
	RunPath    string `json:"run_path"`                         //运行路径
	RunCmd     string `json:"run_cmd"`                          //运行命令
	Num        int    `json:"num"`                              //可重试次数
	Pid        string `json:"pid"`                              //当前运行PID
	PLog       string `json:"p_log"`                            //日志位置
	AutoRun    int    `json:"auto_run" gorm:"default:2"`        //自动运行（单选：1为自动运行  2为不自动运行）
	Remark     string `json:"remark"`                           //备注
	Running    int    `json:"running" gorm:"-"`                 //运行状态 1为未运行 2为运行中
	UpdateUser string `json:"update_user" gorm:"comment:'更新人'"` //更新人
}

// Model GormMode修正
type Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt Time           `json:"createdTime"`
	UpdatedAt Time           `json:"updatedTime"`
	DeletedAt gorm.DeletedAt `json:"deletedTime" gorm:"index"`
}
