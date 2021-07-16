package rpcs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hprose/hprose-golang/v3/io"
	"github.com/hprose/hprose-golang/v3/rpc"
	"github.com/hprose/hprose-golang/v3/rpc/codec/jsonrpc"
	"github.com/hprose/hprose-golang/v3/rpc/core"
	"github.com/hprose/hprose-golang/v3/rpc/plugins/push"
	"github.com/hprose/hprose-golang/v3/rpc/plugins/reverse"
	"github.com/uouuou/ServerManagerSystem/client/cReverse"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/crons"
	"github.com/uouuou/ServerManagerSystem/server/process"
	"github.com/uouuou/ServerManagerSystem/server/registered"
	"github.com/uouuou/ServerManagerSystem/util"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Service 初始化一个service配置
var Service *push.Broker
var Caller *reverse.Caller

//Register 客户端注册服务
type Register struct {
	lock sync.RWMutex
}

// AuthHandler 一个用于权限验证的中间件必须按照如下方式接入NextInvokeHandler用于Handler的使用，还有NextIOHandler 是还没解码序列化的信息
func AuthHandler(c context.Context, name string, args []interface{}, next core.NextInvokeHandler) (response []interface{}, err error) {
	serviceContext := core.GetServiceContext(c)
	// 通过RequestHeaders方法获取string来实现token和userid的获取从而实现校验
	token := serviceContext.RequestHeaders().GetString("token")
	userid := serviceContext.RequestHeaders().GetString("userid")
	timeNow := time.Now().Format("2006-01-02 15")
	tokenNow := mod.Md5V(userid + mid.GetAuth() + timeNow)
	if token == tokenNow {
		return next(c, name, args)
	} else {
		return nil, errors.New("token check failed")
	}

}

// ServiceCompressHandler 使用插件建立的一个gzip加密和压缩方案
func ServiceCompressHandler(ctx context.Context, request []byte, next core.NextIOHandler) (response []byte, err error) {
	req, _ := mid.RpcDecompress(request, nil)
	return mid.RpcCompress(next(ctx, req))
}

// Services 构建一个RPC service
func init() {
	//构建一个rpc服务
	coreService := rpc.NewService()
	service := push.NewBroker(coreService)
	//设置是否启用BUG
	service.Codec = core.NewServiceCodec(core.WithDebug(true))
	//启用日志打印模式
	//service.Use(log.Plugin)
	//中间件验证客户端token
	service.Use(AuthHandler)
	// 插件gzip加密
	service.Use(ServiceCompressHandler)
	// 通过AddInstanceMethods发布接口
	/*  关于类型转换的相关说明以及客户端和服务端的类型转换过程说明
	其它的类型用 interface{} 接收的时候：
	bool -> bool
	nil -> nil
	字符串 -> 字符串
	[]byte -> []byte
	int8, uint8, int16, uint16, int32 -> int
	uint32, int, uint 如果大于 int32 最大值返回为 int64，小于等于 Int32 最大值返回为 int
	int64, uint64, long 返回为 int64
	浮点数返回 -> float64
	时间返回为 -> time.Time
	UUID 返回为 uuid.UUID
	slice 返回为 []interface{}
	map 返回为 map[string]interface{}
	结构体如果注册了返回为注册的结构体指针，未注册返回为 map[string]interface{}
	对于返回 int64 的情况，如果配置了 WithLongType(encoding.LongTypeUint64)，那返回的就是 uint64，如果配置的是：LongTypeBigInt，返回的就是 *big.Int
	对于返回 float64 的情况，如果配置了 WithRealType(encoding.RealTypeFloat32)，那返回的就是 float32，如果配置的是 RealTypeBigFloat 返回的就是 *big.Float
	对于返回 map 类型来说，如果配置了 WithMapType(encoding.MapTypeIIMap)，那返回的就是 map[interface{}]interface{}
	例如：
	client.Codec = rpc.NewClientCodec(
		rpc.WithSimple(true),
		rpc.WithLongType(encoding.LongTypeUint64),
		rpc.WithRealType(encoding.RealTypeFloat64),
		rpc.WithMapType(encoding.MapTypeIIMap),
	)
	服务器端也是可以做同样的配置的
	服务器端加 service.Codec = jsonrpc.NewServiceCodec(nil) jsonrpc的解码器，服务器可以同时支持 hprose 客户端和 jsonrpc 客户端访问
	*/
	service.Codec = jsonrpc.NewServiceCodec(nil)
	service.AddInstanceMethods(&Register{})
	io.RegisterName("process", (*process.Process)(nil))
	io.RegisterName("rpcProcess", (*mod.Process)(nil))
	io.RegisterName("register", (*registered.Register)(nil))
	io.RegisterName("runSql", (*cReverse.RunSql)(nil))
	io.RegisterName("pages", (*mid.Pages)(nil))
	io.RegisterName("myCron", (*crons.MyCron)(nil))
	//通过AddFunction方法发布客户端注册服务
	service.AddFunction(registered.RegisterRpc, "RegisterRpc")
	//注册一个客户端获取本机的定时任务的接口
	service.AddFunction(crons.RpcCronList, "RpcCronList")
	//service.AddFunction(FTime,"Time")
	//设置用户在连接成功后的反映
	socketHandler := rpc.SocketHandler(service)
	Service = service
	mid.Service = service
	Caller = reverse.NewCaller(service.Service)
	mid.Caller = Caller
	// 请求客户端的总体消耗时间超过改时间记为超时
	//Caller.Timeout = time.Second * 30
	// 客户端和服务端的整体通信超时时间
	//Caller.HeartBeat = time.Second * 5
	socketHandler.OnAccept = func(c net.Conn) net.Conn {
		mid.Log().Info(c.RemoteAddr().String() + " -> " + c.LocalAddr().String() + " connected")
		return c
	}
	//设置用户在连接关闭后的反映
	socketHandler.OnClose = func(c net.Conn) {
		mid.Log().Warning(c.RemoteAddr().String() + " -> " + c.LocalAddr().String() + " closed on client")
	}
}

// RunRpc 启动一个TCP RPC服务器
func RunRpc() {
	//使用net监听tcp端口用于数据交互
	server, err := net.Listen("tcp", ":"+strconv.Itoa(util.RpcPort))
	if err != nil {
		mid.Log().Error(err.Error())
		return
	}
	//service绑定监听
	err = Service.Bind(server)
	if err != nil {
		mid.Log().Error(err.Error())
		return
	}
}

// IDL 获取在线名单
func IDL(c *gin.Context) {
	c.JSON(200, gin.H{
		"id": Service.IdList("update"),
	})
}

//func FTime(ctx context.Context,name string)  {
//	serviceContext := rpc.GetServiceContext(ctx)
//	producer := serviceContext.Items().GetInterface("producer").(push.Producer)
//	for  {
//		producer.Push(name+":"+time.Now().String(), "test","233","244","255")
//		time.Sleep(time.Millisecond*300)
//	}
//}

type HaHa struct {
	Name     string `json:"name"`
	Messages string `json:"messages"`
}

func TestHaHa(c *gin.Context) {
	io.RegisterName("haha", (*HaHa)(nil))
	hah := HaHa{
		Name:     "SMS",
		Messages: "Ping",
	}
	clientList, err := ClientList()
	if err != nil {
		mid.DataErr(c, err, "客户端列表获取异常")
		return
	}
	start := Service.Push(hah, "test", clientList...)
	c.JSON(200, gin.H{
		"code":  2000,
		"start": start,
		"data": gin.H{
			"cid": Service.IdList("test"),
			"req": hah,
		},
		"cid":     Service.IdList("test"),
		"list":    clientList,
		"message": "请求成功",
	})
}

func Update(c *gin.Context) {
	data, err := registered.PushUpdate()
	if err != nil {
		mid.DataErr(c, err, "更新信息获取异常")
		return
	}
	io.RegisterName("updateData", (*mod.Update)(nil))
	clientList, err := ClientList()
	if err != nil {
		mid.DataErr(c, err, "客户端列表获取异常")
		return
	}
	start := Service.Push(data, "update", clientList...)
	c.JSON(200, gin.H{
		"start": start,
		"data":  data,
		"ids":   Service.IdList("update"),
		"list":  clientList,
	})
}

// ClientList 获取注册客户端列表
func ClientList() (list []string, err error) {
	var (
		db       = util.GetDB()
		register []registered.Register
		r        registered.Register
	)
	if err = db.Model(&r).Where("deleted_at IS NULL").Find(&register).Error; err != nil {
		return
	}
	for _, registers := range register {
		list = append(list, registers.Userid)
	}
	return
}

// ReverseClient 远程调用客户端方法的结构体
type ReverseClient struct {
	Hello             func(name string) (string, error)                            `name:"Hello"`
	SqlSet            func(sql string) (map[string]interface{}, error)             `name:"SqlSet"`
	SpecifySqlReverse func(runSql cReverse.RunSql) (map[string]interface{}, error) `name:"SpecifySqlReverse"`
}

// ReverseHello 远程调度客户端实现测试
func ReverseHello(c *gin.Context) {
	var (
		results []string
		wg      sync.WaitGroup
	)
	clientList, err := ClientList()
	if err != nil {
		mid.DataErr(c, err, "客户端列表获取异常")
		return
	}
	idList := Caller.IdList()
	for _, i := range clientList {
		for _, s := range idList {
			if reflect.DeepEqual(i, s) {
				wg.Add(1)
				go func(s string) {
					defer wg.Done()
					var (
						err    error
						result string
						proxy  *ReverseClient
					)

					Caller.UseService(&proxy, s)
					result, err = proxy.Hello("world")
					if err != nil {
						mid.Log().Error(mid.RunFuncName() + ":远程调度异常 " + err.Error())
						result = err.Error()
					}
					result = result + ":" + s
					results = append(results, result)
				}(s)
			}
		}
	}
	wg.Wait()
	mid.DataOk(c, results, "请求成功")

}

// ReverseInvokeSql 远程调度客户端实现SQL执行
func ReverseInvokeSql(c *gin.Context) {
	var proxy *ReverseClient
	var re []map[string]interface{}
	sql := c.Query("sql")
	cid := c.Query("id")
	if sql == "" || cid == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	idList := Service.IdList("test")
	for _, s := range idList {
		if s == cid {
			Caller.UseService(&proxy, cid)
			//SELECT * FROM `processes` WHERE `processes`.`deleted_at` IS NULL
			res, err := proxy.SqlSet(sql)
			if err != nil {
				mid.Log().Error(mid.RunFuncName() + ":远程调度异常 " + err.Error())
				mid.DataNot(c, err.Error(), "远程调用失败")
				return
			}
			switch res["data"].(type) {
			case []uint8:
				err = json.Unmarshal(res["data"].([]byte), &re)
				if err != nil {
					mid.DataNot(c, err, "JSON序化化失败")
					return
				}
				res["data"] = re
			}
			c.JSON(http.StatusOK, res)
			return
		}
	}
	mid.DataNot(c, nil, "客户端不在线")
}

// ReverseInvokeEveryOneSql 远程调度客户端实现SQL执行
func ReverseInvokeEveryOneSql(c *gin.Context) {
	var proxy *ReverseClient
	var re []map[string]interface{}
	var run cReverse.RunSql
	err := c.BindJSON(&run)
	if err != nil {
		mid.ClientBreak(c, err, "数据绑定异常")
		return
	}
	if run.Cid == "" || run.SqlRaw == "" || run.SqlConfig.DbName == "" {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	idList := Service.IdList("test")
	for _, s := range idList {
		if s == run.Cid {
			Caller.UseService(&proxy, run.Cid)
			//SELECT * FROM `processes` WHERE `processes`.`deleted_at` IS NULL
			res, err := proxy.SpecifySqlReverse(run)
			if err != nil {
				mid.Log().Error(mid.RunFuncName() + ":远程调度异常 " + err.Error())
				mid.DataNot(c, err.Error(), "远程调用失败")
				return
			}
			switch res["data"].(type) {
			case []uint8:
				err = json.Unmarshal(res["data"].([]byte), &re)
				if err != nil {
					mid.DataNot(c, err, "JSON序化化失败")
					return
				}
				res["data"] = re
			}
			c.JSON(http.StatusOK, res)
			return
		}
	}
	mid.DataNot(c, nil, "客户端不在线")
}

// ReverseInvokeEveryAnySql 远程调度客户端实现SQL执行
func ReverseInvokeEveryAnySql(c *gin.Context) {
	var proxy *ReverseClient
	var re []map[string]interface{}
	var runOnId cReverse.RunSqlOnId
	var run cReverse.RunSql
	var db = util.GetDB()
	err := c.BindJSON(&runOnId)
	if err != nil {
		mid.ClientBreak(c, err, "数据绑定异常")
		return
	}
	if runOnId.Cid == "" || runOnId.SqlRaw == "" || runOnId.SqlId == 0 {
		mid.ClientBreak(c, nil, "参数错误")
		return
	}
	srs := db.Table("sql_registers").Where("id = ? and deleted_at IS NULL", runOnId.SqlId).Find(&run.SqlConfig)
	if srs.Error != nil {
		mid.DataErr(c, err, "查询对应ID异常")
		return
	}
	if srs.RowsAffected <= 0 {
		mid.DataNot(c, nil, "没有查询到对应ID")
		return
	}
	run.Cid = runOnId.Cid
	run.SqlRaw = runOnId.SqlRaw
	idList := Service.IdList("test")
	for _, s := range idList {
		if s == runOnId.Cid {
			Caller.UseService(&proxy, runOnId.Cid)
			//SELECT * FROM `processes` WHERE `processes`.`deleted_at` IS NULL
			res, err := proxy.SpecifySqlReverse(run)
			if err != nil {
				mid.Log().Error(mid.RunFuncName() + ":远程调度异常 " + err.Error())
				mid.DataNot(c, err.Error(), "远程调用失败")
				return
			}
			switch res["data"].(type) {
			case []uint8:
				err = json.Unmarshal(res["data"].([]byte), &re)
				if err != nil {
					mid.DataNot(c, err, "JSON序化化失败")
					return
				}
				res["data"] = re
			}
			c.JSON(http.StatusOK, res)
			return
		}
	}
	mid.DataNot(c, nil, "客户端不在线")
}
