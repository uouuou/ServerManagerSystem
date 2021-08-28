package util

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/models"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Port    int
	RpcPort int
)

type Config struct {
	Setting Setting `json:"setting" yaml:"setting"`
	Sql     Sql     `json:"sql" yaml:"sql"`
}

type Sql struct {
	DbType string `json:"db_type" yaml:"dbType"` // 数据库连接 mysql sqlLite
	DbName string `json:"db_name" yaml:"dbName"` // 如果是mysql dbName
	DbUser string `yaml:"DbUser" json:"dbUser"`  // 如果是mysql dbUser
	DbPass string `yaml:"DbPass" json:"dbPass"`  // 如果是mysql dbPass
	DbHost string `yaml:"DbHost" json:"dbHost"`  // 如果是mysql dbHost
	DbPort int    `yaml:"DbPort" json:"dbPort"`  // 如果是mysql dbPort
}

type Setting struct {
	Port    int    `json:"port" yaml:"port"`        //程序运行端口
	RunDir  string `json:"run_dir" yaml:"runDir"`   //程序运行路径
	RedType string `json:"red_type" yaml:"redType"` // sftp读取文本的后缀名使用 | 隔开
	RpcPort int    `json:"rpc_port" yaml:"rpcPort"` //内部RPC通信端口默认是 8001
	Auth    string `json:"auth" yaml:"auth"`        //服务端授权码（客户端设置一致即可连接）
}

type Client struct {
	Port       int    `json:"port" yaml:"port"`             //程序运行端口
	RunDir     string `json:"run_dir" yaml:"runDir"`        //程序运行路径
	Server     string `json:"server" yaml:"server"`         //远端服务端RPC地址
	Userid     string `json:"userid" yaml:"userid"`         //远程注册ID
	ServerHttp string `json:"serverHttp" yaml:"serverHttp"` //远程服务端HTTP地址
	Auth       string `json:"auth" yaml:"auth"`             //服务端授权码（服务端设置一致即可）
}

var DB *gorm.DB

func init() {
	var sqlErr error
	system := runtime.GOOS
	switch system {
	case "linux":
		{
			//检查设备是否可写入，若不可写入则重启设备
			touch := exec.Command("touch", "-a", "/opt/readonly_test")
			reboot := exec.Command("reboot")
			err := touch.Run()
			if err != nil {
				mid.Log().Info(fmt.Sprintf("err:%v", err))
				err = reboot.Run()
				if err != nil {
					mid.Log().Info(fmt.Sprintf("err:%v", err))
				}
			} else {
				rm := exec.Command("rm", "-rf", "/opt/readonly_test")
				err = rm.Run()
				if err != nil {
					mid.Log().Info(fmt.Sprintf("err:%v", err))
				}
			}
		}
	}
	// 从本地读取环境变量+
	loadErr := godotenv.Load()
	if loadErr != nil {
		gin.SetMode(gin.ReleaseMode)
	}
	if mid.GetMode() == "client" {
		// 获取配置文件信息
		var conf Client
		configFile := mid.Dir + "/config/client.yaml"
		config, err := os.ReadFile(configFile)
		if err != nil {
			_, err = os.Create(configFile)
			if err != nil {
				mid.Log().Error(err.Error())
			} else {
				config, err = os.ReadFile(configFile)
				if err != nil {
					mid.Log().Error(err.Error())
				}
			}
		}
		err = yaml.Unmarshal(config, &conf)
		if err != nil {
			mid.Log().Error(err.Error())
		}
		Port = conf.Port
		if Port == 0 || conf.Server == "" || conf.Userid == "" || conf.ServerHttp == "" || conf.Auth == "" {
			conf.Port = 8002
			conf.Server = "tcp://127.0.0.1:8001"
			conf.ServerHttp = "http://127.0.0.1:8000"
			conf.Auth = "ServerManagerSystem2021"
			conf.Userid = mid.CId()

			data, err := yaml.Marshal(conf)
			if err != nil {
				mid.Log().Error(err.Error())
			}
			err = os.WriteFile(configFile, data, 0777)
			if err != nil {
				mid.Log().Error(err.Error())
			}
		}
		mid.CUId = conf.Userid
		mid.Server = conf.Server
		mid.ServerHttp = conf.ServerHttp
		mid.Auth = conf.Auth
		dsn := mid.Dir + "/config/client.db"
		if !models.FileExist(dsn) {
			configDir := mid.Dir + "/config"
			_, err = os.Stat(configDir)
			if err != nil {
				if os.IsNotExist(err) {
					err := os.MkdirAll(configDir, os.ModePerm)
					if err != nil {
						mid.Log().Info(fmt.Sprintf("err:%v", err))
					}
					return
				}
				return
			}
			// 创建一个数据库文件
			file, err := os.Create(dsn)
			if err != nil {
				mid.Log().Info(fmt.Sprintf("err:%v", err))
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
		}
		DB, sqlErr = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	} else {
		// 获取配置文件信息
		var conf Config
		configFile := mid.Dir + "/config/config.yaml"
		config, err := os.ReadFile(configFile)
		if err != nil {
			_, err = os.Create(configFile)
			if err != nil {
				mid.Log().Error(err.Error())
			} else {
				config, err = os.ReadFile(configFile)
				if err != nil {
					mid.Log().Error(err.Error())
				}
			}
		}
		err = yaml.Unmarshal(config, &conf)
		if err != nil {
			mid.Log().Error(err.Error())
		}
		Port = conf.Setting.Port
		RpcPort = conf.Setting.RpcPort
		if Port == 0 || conf.Sql.DbType == "" || conf.Setting.RedType == "" || conf.Setting.RpcPort == 0 || conf.Setting.Auth == "" {
			conf.Setting.Port = 8000
			conf.Setting.RpcPort = 8001
			conf.Setting.Auth = "ServerManagerSystem2021"
			conf.Sql.DbType = "sqllite"
			conf.Setting.RedType = ".txt|.sh|.log|.config|.ini|.in|.md"
			data, err := yaml.Marshal(conf)
			if err != nil {
				mid.Log().Error(err.Error())
			}
			err = os.WriteFile(configFile, data, 0777)
			if err != nil {
				mid.Log().Error(err.Error())
			}
		}
		// 设置中间件读取一些需要的配置信息
		mid.ReadType = conf.Setting.RedType
		mid.Auth = conf.Setting.Auth
		// 设置数据库
		switch conf.Sql.DbType {
		case "mysql":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Sql.DbUser, conf.Sql.DbPass, conf.Sql.DbHost, conf.Sql.DbPort, conf.Sql.DbName)
			DB, sqlErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger:                 logger.Default.LogMode(logger.Error),
				PrepareStmt:            true,
				SkipDefaultTransaction: true,
			})
		case "sqllite":
			dsn := mid.Dir + "/config/server.db"
			if !models.FileExist(dsn) {
				configDir := mid.Dir + "/config"
				_, err := os.Stat(configDir)
				if err != nil {
					if os.IsNotExist(err) {
						err := os.MkdirAll(configDir, os.ModePerm)
						if err != nil {
							mid.Log().Info(fmt.Sprintf("err:%v", err))
						}
						return
					}
					return
				}
				// 创建一个数据库文件
				file, err := os.Create(dsn)
				if err != nil {
					mid.Log().Info(fmt.Sprintf("err:%v", err))
				}
				defer func(file *os.File) {
					_ = file.Close()
				}(file)
			}
			DB, sqlErr = gorm.Open(sqlite.Open(dsn), &gorm.Config{
				Logger:                 logger.Default.LogMode(logger.Error),
				PrepareStmt:            true,
				SkipDefaultTransaction: true,
			})
		}
	}

	if sqlErr != nil {
		mid.Log().Error(fmt.Sprintf("Sql connect error %v\n", sqlErr))
	}
	if DB.Error != nil {
		mid.Log().Error(fmt.Sprintf("database error %v\n", DB.Error))
	}
	sqlDB, err := DB.DB()
	if err != nil {
		mid.Log().Error(fmt.Sprintf("sql err:%v\n", err))
	}
	// 连接池最多同时打开的连接数
	sqlDB.SetMaxOpenConns(200)
	//连接池里最大空闲连接数
	sqlDB.SetMaxIdleConns(50)
	//连接池里面的连接最大存活时长
	sqlDB.SetConnMaxLifetime(time.Hour * 2)
	//连接池里面的连接最大空闲时长
	sqlDB.SetConnMaxIdleTime(time.Hour)
	if mid.GetMode() == "server" && !mid.Version {
		models.NewRoutine(RouterInit)
	}
}

func GetDB() *gorm.DB {
	return DB
}

// RouterInit 为路由获取一个菜单信息
func RouterInit() {
	var menu []mid.Menu
	if err := DB.Model(&menu).Where("authority = 1 and deleted_at IS NULL").Find(&menu).Error; err != nil {
		mid.Log().Error(err.Error())
	}
	mid.MenuList = menu
}
