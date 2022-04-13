package rpcc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hprose/hprose-golang/v3/io"
	"github.com/hprose/hprose-golang/v3/rpc"
	"github.com/hprose/hprose-golang/v3/rpc/core"
	"github.com/hprose/hprose-golang/v3/rpc/plugins/push"
	"github.com/hprose/hprose-golang/v3/rpc/plugins/reverse"
	"github.com/uouuou/ServerManagerSystem/client/cCron"
	"github.com/uouuou/ServerManagerSystem/client/cProcess"
	"github.com/uouuou/ServerManagerSystem/client/cReverse"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	con "github.com/uouuou/ServerManagerSystem/server"
	"github.com/uouuou/ServerManagerSystem/server/registered"
	"github.com/uouuou/ServerManagerSystem/server/rpcs"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Prosumer 初始化构建一个用于push和订阅的Prosumer
var Prosumer = ProsumerFun()

// Client 初始化构建一个用于连接的连接器配置
var Client = RpcClient()

//Provider 注册远程方法发布
var Provider = MakeProvider()

// UpdateStatus 更新进行中的状态
var UpdateStatus bool

type Update struct {
	Register func(r registered.Register) (mid.ResultBody, error) `name:"RegisterRpc"`
	Login    func(userid string, token string) (bool, error)     `name:"Login"`
}

type Process struct {
	AddRpcProcess   func(process mod.Process) (mid.ResultBody, error)                             `name:"AddRpcProcess"`
	ProcessRpcList  func(userid string) (mid.ResultBody, error)                                   `name:"ProcessRpcList"`
	ProcessRpcLists func(userid string, page int, pageSize int) ([]mod.Process, mid.Pages, error) `name:"ProcessRpcLists"`
}

// RpcC 建立客户端与服务端的通信
func RpcC() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("panic error.")
			}
		}()
		// 注册客户端
		RegisterFrpc(Client)
	}()
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("panic error.")
			}
		}()
		// 注册一个定时任务接口
		cCron.AutoCronRun(Client)
	}()
	mod.NewRoutine(TimeTest)
	mod.NewRoutine(UpdateClient)
	mod.NewRoutine(Reverse)
}

// RpcClient 建立RPC链接
func RpcClient() *rpc.Client {
	client := rpc.NewClient(mid.GetServer())
	//client.Use(log.Plugin)
	//验证客户端合法性的插件
	client.Use(AuthHandler)
	//插件gzip加密
	client.Use(ClientCompressHandler)
	socketTransport := rpc.SocketTransport(client)
	socketTransport.OnConnect = func(c net.Conn) net.Conn {
		mid.Log.Info(c.LocalAddr().String() + " -> " + c.RemoteAddr().String() + " connected")
		return c
	}
	socketTransport.OnClose = func(c net.Conn) {
		mid.Log.Warning(c.LocalAddr().String() + " -> " + c.RemoteAddr().String() + " closed on client")
	}
	return client
}

// AuthHandler 一个用户客户端验证的插件 NextInvokeHandler用于Handler的使用，还有NextIOHandler 是还没解码序列化的信息
func AuthHandler(ctx context.Context, name string, args []interface{}, next core.NextInvokeHandler) (result []interface{}, err error) {
	timestamp := time.Now().Unix()
	token := mod.Md5V(mid.GetCUId() + mid.GetAuth() + strconv.FormatInt(timestamp, 10))
	headers := core.GetClientContext(ctx).RequestHeaders()
	headers.Set("token", token)
	headers.Set("userid", mid.GetCUId())
	headers.Set("timestamp", timestamp)
	return next(ctx, name, args)
}

// ClientCompressHandler 使用插件构建的一个gzip数据压缩
func ClientCompressHandler(ctx context.Context, request []byte, next core.NextIOHandler) (response []byte, err error) {
	req, _ := mid.RpcCompress(request, nil)
	return mid.RpcDecompress(next(ctx, req))
}

// ProsumerFun 构建一个push.Prosumer方法
func ProsumerFun() *push.Prosumer {
	//client.Use(log.Plugin)
	prosumer := push.NewProsumer(Client, mid.GetCUId())
	prosumer.OnSubscribe = func(topic string) {
		mid.Log.Info(topic + " is subscribed.")
	}
	prosumer.OnUnsubscribe = func(topic string) {
		mid.Log.Info(topic + " is unsubscribed.")
	}
	prosumer.RetryInterval = time.Second * 2
	return prosumer
}

// MakeProvider 创建一个客户端方法提供者
func MakeProvider() *reverse.Provider {
	provider := reverse.NewProvider(Client, mid.GetCUId())
	provider.Debug = true
	return provider
}

// Reverse 反向调度方法发布
func Reverse() {
	io.RegisterName("runSql", (*cReverse.RunSql)(nil))
	io.RegisterName("rpcProcess", (*mod.Process)(nil))
	io.RegisterName("pages", (*mid.Pages)(nil))
	Provider.AddFunction(cReverse.Hello, "Hello")
	Provider.AddFunction(cReverse.SqlReverse, "SqlSet")
	Provider.AddFunction(cReverse.SpecifySqlReverse, "SpecifySqlReverse")
	Provider.AddInstanceMethods(&cReverse.CProcess{})
	mod.NewRoutine(Provider.Listen)
}

// RegisterFrpc 用于客户端与服务端的注册和更新状态以及发现
func RegisterFrpc(client *rpc.Client) bool {
	var stub *Update
	var r registered.Register
	//使用接口传递参数需要先注册对应的结构体，下面是针对该结构体进行注册
	io.RegisterName("register", (*registered.Register)(nil))
	client.UseService(&stub)
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		r.Userid = mid.GetCUId()
		r.ClientIp = mod.GetLocalIP()
		r.Version = mod.AppVersion
		if r.Version == "" {
			r.Version = gin.Mode()
		}
		frpVersion := mod.FrpVersion()
		if frpVersion == "" {
			frpVersion = "not install"
		}
		r.FrpVersion = frpVersion
		npsVersion := mod.NpsVersion()
		if npsVersion == "" {
			npsVersion = "version: not install"
		}
		npsVersionNow := strings.Split(npsVersion, ": ")
		r.NpsVersion = npsVersionNow[1]
		register, err := stub.Register(r)
		if err != nil {
			mid.Log.Error(mid.RunFuncName() + ": RPC调度异常：" + err.Error())
			continue
		}
		switch register.Code {
		case 2000:
			rs := register.Data.(*registered.Register)
			mid.NatAuth = rs.NatAuth
			mid.CronAuth = rs.CronAuth
			mid.Log.Info(register.Message)
		case 2001:
			rs := register.Data.(*registered.Register)
			mid.NatAuth = rs.NatAuth
			mid.CronAuth = rs.CronAuth
			//将远程的配置文件写入到本地(文件和变量)
			if !reflect.DeepEqual(mid.GetNPSConfig(), rs.NpsConfig) {
				if rs.NpsConfig != "null" {
					var p cProcess.Process
					mid.NPSConfig = rs.NpsConfig
					plist, _ := p.ProcessRpcLists(1, -1)
					for _, process := range plist {
						if process.Name == "NPC" {
							process.RunCmd = mid.Dir + "/config/nps/npc -server=" + mid.GetNPSConfig()
							ps := mod.Process{
								Model: mod.Model{
									ID: process.ID,
								},
								Name:    process.Name,
								RunPath: process.RunPath,
								RunCmd:  mid.Dir + "/config/nps/npc -server=" + mid.GetNPSConfig(),
								Num:     process.Num,
								AutoRun: process.AutoRun,
								Remark:  process.Remark,
							}
							res := p.EditFunProcess(ps)
							if res.Code == 2000 {
								resDate := p.OffRpcProcess(process)
								if resDate.Code == 2000 {
									mid.Log.Info(mid.RunFuncName() + ":" + resDate.Message)
								} else {
									mid.Log.Error(mid.RunFuncName() + ":" + resDate.Message)
								}
								mid.Log.Info(mid.RunFuncName() + ":" + res.Message)
							} else {
								mid.Log.Error(mid.RunFuncName() + ":" + res.Message)
							}
						}
					}
					err = os.WriteFile(mid.Dir+"/config/nps/npc.conf", []byte(rs.NpsConfig), 0755)
					if err != nil {
						mid.Log.Error(mid.RunFuncName() + " :数据写入失败 " + err.Error())
					}
					mid.Log.Info("NPC配置文件已更新")
				}
			}
			if !reflect.DeepEqual(mid.GetFRTPConfig(), rs.FrpConfig) {
				if rs.FrpConfig != "null" {
					var p cProcess.Process
					mid.FRPConfig = rs.FrpConfig
					plist, _ := p.ProcessRpcLists(1, -1)
					for _, process := range plist {
						if process.Name == "FRPC" {
							p.OffRpcProcess(process)
						}
					}
					err = os.WriteFile(mid.Dir+"/config/frp/frpc.ini", []byte(rs.FrpConfig), 0755)
					if err != nil {
						mid.Log.Error(mid.RunFuncName() + " :数据写入失败 " + err.Error())
					}
					mid.Log.Info("FRP配置文件已更新")
				}
			}
		case 2003:
			mid.Log.Error(mid.RunFuncName() + ":" + register.Message)
		}
	}
	return false
}

func TimeTest() {
	io.RegisterName("haha", (*rpcs.HaHa)(nil))
	_, err := Prosumer.Subscribe("test", func(data *rpcs.HaHa, from string) {
		if data.Messages == "Ping" {
			res := rpcs.HaHa{
				Name:     mid.GetCUId() + "data set  is ok",
				Messages: "Pang",
			}
			m, err := Prosumer.Push(res, "test", mid.GetCUId())
			if err != nil {
				mid.Log.Error(mid.RunFuncName() + " :err " + err.Error())
			}
			mid.Log.Infof("push info %v", m)
		}
	})
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + err.Error())
	}
	/* 远程调用方法
	result, err := client.Invoke("Time", []interface{}{"too difficult"})
	if err != nil {
		return
	}
	fmt.Println(result)
	*/
}

func UpdateClient() {
	io.RegisterName("updateData", (*mod.Update)(nil))
	_, err := Prosumer.Subscribe("update", func(data *mod.Update, from string) {
		if UpdateStatus {
			mid.Log.Warning("有更新任务正在继续，请稍后重试......")
			return
		}
		if gin.Mode() == "release" {
			var cmdStart *exec.Cmd
			UpdateStatus = true
			architecture, system := mid.GetArchitecture()
			cmd := exec.Command("chmod", "-R", "777", mid.Dir+"/config/"+mid.MainName+"_tmp")
			cmdMv := exec.Command("mv", mid.Dir+"/config/"+mid.MainName+"_tmp", mid.Dir+"/"+mid.MainName)
			switch architecture {
			case "arm64":
				if !con.Down(mid.GetServerHttp()+data.UrlArm, mid.Dir+"/config/"+mid.MainName+"_tmp") {
					mid.Log.Error(mid.RunFuncName() + ":下载更新文件出现异常")
					UpdateStatus = false
					return
				}
				cmdStart = exec.Command("systemctl", "restart", "serverManager_"+mid.Mode)
			case "amd64":
				if !con.Down(mid.GetServerHttp()+data.UrlLinux, mid.Dir+"/config/"+mid.MainName+"_tmp") {
					mid.Log.Error(mid.RunFuncName() + ":下载更新文件出现异常")
					UpdateStatus = false
					return
				}
				cmdStart = exec.Command("systemctl", "restart", "serverManager_"+mid.Mode)
			default:
				cmdStart = exec.Command("service", "serverManager_"+mid.Mode, "restart")
			}
			switch system {
			case "windows":
				{
					mid.Log.Info(mid.RunFuncName() + ":暂时不支持WIN系统升级")
					UpdateStatus = false
					return
				}
			case "linux":
				{
					err := cmd.Run()
					if err != nil {
						mid.Log.Info(mid.RunFuncName() + ":CMD执行错误 " + err.Error())
						UpdateStatus = false
						return
					}

					err = cmdMv.Run()
					if err != nil {
						mid.Log.Info(mid.RunFuncName() + ":CMD执行错误 " + err.Error())
						UpdateStatus = false
						return
					}

					mid.Log.Info(data.Version + ":更新已经完成,即将重启")

					err = cmdStart.Run()
					if err != nil {
						mid.Log.Info(mid.RunFuncName() + ":CMD执行错误 " + err.Error())
						UpdateStatus = false
						return
					}

				}
			default:
				mid.Log.Info(mid.RunFuncName() + ":暂时不支持" + system + "系统升级")
				UpdateStatus = false
				return
			}
			UpdateStatus = false
		} else {
			mid.Log.Infof("更新版本：%v  linux版本地址：%v arm版本地址：%v 备注：%v", data.Version, data.UrlLinux, data.UrlArm, data.Remark)
		}
	})
	if err != nil {
		mid.Log.Error(mid.RunFuncName() + err.Error())
	}
}
