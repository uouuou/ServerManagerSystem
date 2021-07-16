package cmd

import (
	"github.com/kardianos/service"
	"github.com/uouuou/ServerManagerSystem/client"
	"github.com/uouuou/ServerManagerSystem/client/cNat"
	"github.com/uouuou/ServerManagerSystem/client/rpcc"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/middleware/convert"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/routes"
	con "github.com/uouuou/ServerManagerSystem/server"
	"github.com/uouuou/ServerManagerSystem/server/clash"
	"github.com/uouuou/ServerManagerSystem/server/process"
	"github.com/uouuou/ServerManagerSystem/server/rpcs"
	"github.com/uouuou/ServerManagerSystem/server/system"
	"log"
	"os"
)

type Program struct{}

func (p *Program) Start(service.Service) error {
	log.Println("开始服务...")
	go p.run()
	return nil
}
func (p *Program) Stop(service.Service) error {
	log.Println("停止服务...")
	return nil
}
func (p *Program) run() {
	mod.SystemVersion()
	//根据传入参数注册路由并设置基础信息
	switch mid.Mode {
	case "client":
		//安装NPS
		mod.NewRoutine(cNat.NpsInstall)
		//安装FRP
		mod.NewRoutine(cNat.FrpInstall)
		//注册数据库构建
		mod.NewRoutine(client.IntSqlStart)
		//注册一个定时任务每隔一分钟检查自动启动的程序是否启动
		mod.NewRoutine(process.AutoRun)
		//启动一个RPC连接：通过RPC与服务端做数据交互（更新或者服务端的消息广播）
		mod.NewRoutine(rpcc.RpcC)
		//注册路由
		routes.ClientRouter()

	case "server":
		//初始化操作系统程序
		go func() {
			m := EnvTesting()
			mid.Log().Info(convert.ToString(m["unbound"] + "  " + m["clash"] + "  " + m["nftables"]))
		}()
		//使用goroutine防止单独进程崩溃的情况下运行
		mod.NewRoutine(con.IntSqlStart)    //启动时自动处理SqLite
		mod.NewRoutine(system.CollectTask) //统计程序启动后的网卡上下载情况
		mod.NewRoutine(clash.ReadConfig)   //读取一个clash配置文件
		//mod.NewRoutine(clash.Runs)       //启动一个clash
		mod.NewRoutine(EnvNftables)     //启动一个环境监测程序
		mod.NewRoutine(process.AutoRun) //注册一个定时任务每隔一分钟检查自动启动的程序是否启动
		mod.NewRoutine(rpcs.RunRpc)     //启动rpc服务器
		//安装NPS
		mod.NewRoutine(cNat.NpsInstall)
		//安装FRP
		mod.NewRoutine(cNat.FrpInstall)
		//注册路由
		routes.ServerRouter()

	}
}

func Install() {
	//服务初始化
	//服务的配置信息
	cfg := &service.Config{
		Name:        "serverManagerSystem_" + mid.Mode,
		DisplayName: "ServerManager System",
		Description: "This is a ServerManager System.",
		Arguments:   []string{"-m", mid.Mode},
	}
	// Interface 接口
	prg := &Program{}
	// 构建服务对象
	s, err := service.New(prg, cfg)
	if err != nil {
		log.Fatal(err)
	}
	// logger 用于记录系统日志
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	for i, arg := range os.Args {
		if arg == "-m" {
			mid.Mode = os.Args[i+1]
			if len(os.Args) == i+3 {
				if os.Args[i+2] == "install" {
					err = service.Control(s, "install")
					if err != nil {
						log.Fatal(err)
						return
					}
					log.Fatal(cfg.Name + "服务注册成功：请使用systemctl start " + cfg.Name + "启动程序或通过status查看状态")
					return
				} else if os.Args[i+2] == "uninstall" {
					err = service.Control(s, "uninstall")
					if err != nil {
						log.Fatal(err)
						return
					}
					log.Fatal(cfg.Name + "服务删除成功")
					return
				} else {
					log.Fatal("注册服务或卸载服务请使用sm -m client install | sm -m client uninstall")
				}
			}
		}
	}
	if len(os.Args) == 2 { //如果有命令则执行
		if mid.Version {
			mod.SystemVersion()
			return
		}
		if mid.Readable {
			mod.TestReadable()
			return
		}
		switch os.Args[1] {
		case "start":
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
		case "stop":
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
		case "install":
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Fatal(cfg.Name + "服务注册成功：请使用systemctl start " + cfg.Name + "启动程序或通过status查看状态")
		case "uninstall":
			err = service.Control(s, os.Args[1])
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Fatal(cfg.Name + "服务删除成功")
		}
	} else { //否则说明是方法启动了
		err = s.Run()
		if err != nil {
			_ = logger.Error(err)
		}
	}
}
