package routes

import (
	"context"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	con "github.com/uouuou/ServerManagerSystem/server"
	"github.com/uouuou/ServerManagerSystem/server/clash"
	"github.com/uouuou/ServerManagerSystem/server/crons"
	"github.com/uouuou/ServerManagerSystem/server/logview"
	"github.com/uouuou/ServerManagerSystem/server/menu"
	"github.com/uouuou/ServerManagerSystem/server/nat"
	"github.com/uouuou/ServerManagerSystem/server/net/firewall"
	"github.com/uouuou/ServerManagerSystem/server/net/webshell"
	"github.com/uouuou/ServerManagerSystem/server/process"
	"github.com/uouuou/ServerManagerSystem/server/registered"
	"github.com/uouuou/ServerManagerSystem/server/rpcs"
	"github.com/uouuou/ServerManagerSystem/server/system"
	"github.com/uouuou/ServerManagerSystem/server/upload"
	"github.com/uouuou/ServerManagerSystem/server/user"
	"github.com/uouuou/ServerManagerSystem/util"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// ServerRouter 服务端路由
func ServerRouter() {
	gin.ForceConsoleColor()
	r := gin.New()
	r.NoRoute(mid.HandleNotFound)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(mid.Cors())
	t := template.Must(template.New("").ParseFS(mid.FS, "web/*.html"))
	r.SetHTMLTemplate(t)
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	static, _ := fs.Sub(mid.FS, "web/static")
	r.StaticFS("/upload", http.Dir(mid.Dir+"/upload"))
	r.StaticFS("/static", http.FS(static))
	r.GET("favicon.ico", func(c *gin.Context) {
		file, _ := mid.FS.ReadFile("web/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/404", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	for _, list := range mid.MenuList {
		r.GET(list.Url, func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
	}
	v1 := r.Group("api/v1")
	v1.Use(mid.JWTAuthMiddleware())
	{
		// 系统设置相关接口
		setting := v1.Group("setting")
		setting.Use(mid.AuthCheckMiddleware)
		{
			//菜单管理接口组
			menus := setting.Group("menu")
			{
				//获取所有菜单详情涵盖被屏蔽
				menus.GET("list", menu.GetMenuLists)
				//菜单新增接口
				menus.POST("add", menu.AddMenu)
				//修改菜单信息
				menus.PUT("edit", menu.EditMenu)
				//删除菜单
				menus.DELETE("del", menu.DelMenu)
			}
			//文件管理接口组
			uploads := setting.Group("file")
			{
				//上传文件的管理
				up := upload.Upload{}
				uploads.DELETE("del", up.Del)
				uploads.GET("list", up.List)
			}
			//用户设置接口组
			users := setting.Group("users")
			{
				//获取所有用户详情
				users.GET("list", user.ListUser)
				//用户名新增接口
				users.POST("add", user.AddUser)
				//修改用户信息
				users.PUT("edit", user.EditUser)
				//删除用户
				users.DELETE("del", user.DelUser)
			}
			//角色管理接口组
			role := setting.Group("role")
			{
				// 角色相关接口
				roles := user.Role{}
				role.GET("list", roles.List)
				role.POST("add", roles.Add)
				role.PUT("edit", roles.Edit)
				role.DELETE("del", roles.Del)

			}

		}
		//行为管理接口组
		action := v1.Group("action")
		action.Use(mid.AuthCheckMiddleware)
		{
			//接入管理接口组
			register := action.Group("register")
			{
				//客户端注册的相关接口
				cl := registered.Register{}
				cp := process.CProcessFun{}
				//注册客户端
				register.POST("register", cl.Register)
				//查询客户端列表
				register.GET("cl_list", cl.List)
				//设置客户端允许接入服务器NPS和FRP
				register.POST("cl_set", cl.Set)
				//删除某一个客户端的注册
				register.DELETE("cl_del", cl.Del)
				//查看客户端NPS配置
				register.GET("cl_nps_conf", cl.NpsConf)
				//查看客户端FRP配置
				register.GET("cl_frp_conf", cl.FrpConf)
				//获取客户端守护信息
				register.GET("cp_list", cp.List)
				//新增客户端守护
				register.POST("cp_add", cp.Add)
				//修改客户端守护
				register.PUT("cp_edit", cp.Edit)
				//删除客户端守护
				register.DELETE("cp_del", cp.Del)
				//启动客户端保持
				register.POST("cp_run", cp.Run)
				//关闭客户端保持
				register.POST("cp_off", cp.Off)
			}
			//SQL资产接口组
			sql := action.Group("sql")
			{
				// 和远程调度SQL资产相关的设置
				sqlRegister := registered.SqlRegister{}
				sql.GET("sqlList", sqlRegister.SqlList)
				sql.DELETE("delSql", sqlRegister.DelSql)
				sql.POST("addSql", sqlRegister.AddSql)
				sql.PUT("edit", sqlRegister.EditSql)
				sql.POST("sql_any", rpcs.ReverseInvokeEveryAnySql)
				sql.GET("cid", rpcs.TestHaHa)

			}
			//更新管理接口组
			update := action.Group("update")
			{
				//设置程序更新版本号
				update.POST("set_update_version", registered.SetUpdateVersion)
				//删除更新程序版本号URL的设置
				update.DELETE("del_update_version", registered.DelUpdateVersion)
				//客户端版本更新信息获取
				update.GET("get_update_version", registered.GetUpdateVersion)
			}

		}
		//代理设置接口组
		proxy := v1.Group("clash")
		proxy.Use(mid.AuthCheckMiddleware)
		{
			//订阅服务接口组
			subs := proxy.Group("sub")
			{
				//订阅相关接口
				sub := clash.Sub{}
				subSet := clash.SubSet{}
				//订阅列表
				subs.GET("subList", sub.GetSubList)
				//新增订阅
				subs.POST("addSub", sub.AddSub)
				//修改订阅
				subs.PUT("editSub", sub.EditSub)
				//删除订阅
				subs.DELETE("delSub", sub.DelSub)
				//立即订阅
				subs.GET("subNow", subSet.SubNow)
				//订阅规则设置
				subs.POST("subSet", subSet.SubSet)
				//获取订阅设置
				subs.GET("getSubSet", subSet.GetSubSet)
			}
			// 获取clash设置
			proxy.GET("clashInfo", clash.GetClashInfo)
			//获取clash版本信息
			proxy.GET("onlineClash", clash.OnlineClashVersion)
			//更新clash版本
			proxy.GET("updateClash", clash.UpdateClashNew)
		}
		//网络安全接口组
		net := v1.Group("net")
		net.Use(mid.AuthCheckMiddleware)
		{
			//防火墙接口组
			firewallGroup := net.Group("firewall")
			{
				//服务器防火墙设置
				firewalls := firewall.Firewall{}
				//防火墙列表
				firewallGroup.GET("firewallList", firewalls.FirewallList)
				//新增防火墙
				firewallGroup.POST("addFirewall", firewalls.AddFirewall)
				//删除防火墙
				firewallGroup.DELETE("delFirewall", firewalls.DelFirewall)
			}
			//WebShell接口组
			shell := net.Group("shell")
			{
				//注册一个webSsh相关的接口
				serverInfo := webshell.ServerInfo{}
				//shell服务器列表
				shell.GET("shellList", serverInfo.GetShellList)
				//新增shell服务器
				shell.POST("addShell", serverInfo.AddShell)
				//删除shell服务器
				shell.DELETE("delShell", serverInfo.DelShell)
				//修改shell服务器
				shell.PUT("editShell", serverInfo.EditShell)
				// 查询对应服务器的系统信息
				shell.POST("systemInfo", webshell.ServerSystemInfo)
			}
			//sftp相关接口(接口组)
			vs := net.Group("sftp")
			{
				//ls 接口
				vs.GET("ls", webshell.SftpLs)
				//rm 接口
				vs.DELETE("rm", webshell.SftpRm)
				//创建新目录接口
				vs.GET("mkdir", webshell.SftpMkdir)
				//重命名接口
				vs.GET("rename", webshell.SftpRenameEndpoint)
				//通过sftp上传文件接口
				vs.POST("upload", webshell.SftPUpload)
				//查看文件接口
				vs.GET("cat", webshell.SftpCat)
			}
			//WebShell接口组
			cron := net.Group("cron")
			{
				myCron := crons.MyCron{}
				//Cron列表
				cron.GET("list", myCron.List)
				//新增Cron
				cron.POST("add", myCron.Add)
				//修改Cron
				cron.PUT("edit", myCron.Edit)
				//删除Cron
				cron.DELETE("del", myCron.Del)
			}
		}
		//进程保持接口组
		processGroup := v1.Group("process")
		processGroup.Use(mid.AuthCheckMiddleware)
		{
			//进程保持相关接口
			manage := processGroup.Group("manage")
			{
				//进程守护
				pro := process.Process{}
				//守护进程列表
				manage.GET("process", pro.ProcessList)
				//新增守护
				manage.POST("process", pro.AddProcess)
				//修改守护
				manage.PUT("process", pro.EditProcess)
				//删除守护
				manage.DELETE("process", pro.DelProcess)
				//运行守护
				manage.POST("runProcess", pro.RunProcess)
				//关闭守护程序
				manage.POST("offProcess", pro.OffProcess)
			}

		}
		natGroup := v1.Group("nat")
		natGroup.Use(mid.AuthCheckMiddleware)
		{
			//frp相关接口
			frp := natGroup.Group("frp")
			{
				//查看github在线版本信息
				frp.GET("online", nat.OnlineFrpVersion)
				//更新github在线版本
				frp.GET("update", nat.UpdateFrpNew)
				//获取frps安装情况及配置信息
				frp.GET("get_conf", nat.GetFrpConfig)
				//设置frps配置信息
				frp.POST("set_conf", nat.SetFrpConfig)
			}
			//nps相关接口
			nps := natGroup.Group("nps")
			{
				//查看github在线版本信息
				nps.GET("online", nat.OnlineNpsVersion)
				//更新github在线版本
				nps.GET("update", nat.UpdateNpsNew)
				//获取nps安装情况及配置信息
				nps.GET("get_conf", nat.GetNpsConfig)
				//设置nps配置信息
				nps.POST("set_conf", nat.SetNpsConfig)
			}
			//远程透传功能相关接口
			//获取nps和frp的可用性信息
			natGroup.GET("get_conf", nat.GetNatConf)
			//设置nps和frp的相关配置及是否开启等
			natGroup.POST("set_conf", nat.SetNatConf)

		}
		//获取当前交互的用户信息
		v1.GET("userinfo", mid.MyClaims{}.HomeHandler)
		//修改密码
		v1.PUT("changePass", user.ChangePassword)
		//退出登录
		v1.GET("logout", con.Logout)
		//获取菜单(根据权限)
		v1.GET("menu", menu.GetMenuList)
		// 上传接口（公共接口）
		up := upload.Upload{}
		v1.POST("upload", up.FilesUpload)
		//服务器基础环境信息接口
		v1.GET("systemInfo", system.GetServerInfo)
	}
	notAuth := r.Group("api/open")
	{
		//系统登录接口
		notAuth.POST("login", con.Login)
		//获取程序授权token
		notAuth.GET("gettoken", con.GetToken)
		//注册一个webSsh相关的接口
		serverInfo := webshell.ServerInfo{}
		notAuth.GET("xterm", serverInfo.Xterm)
		//ssh接口
		notAuth.GET("ws", serverInfo.Ws)
		//日志动态显示
		notAuth.GET("ws_log", logview.LogWs)
		//sftp下载文件或压缩后下载文件夹
		notAuth.GET("down", webshell.SftpDownload)
		//rpc相关测试接口
		//notAuth.GET("ids", rpcs.IDL)
		//notAuth.GET("push", rpcs.TestHaHa)
		//notAuth.GET("update", rpcs.Update)
		//notAuth.GET("re", rpcs.ReverseHello)
		//notAuth.GET("sql", rpcs.ReverseInvokeSql)
		//notAuth.POST("sql_one", rpcs.ReverseInvokeEveryOneSql)
		//notAuth.POST("sql_any", rpcs.ReverseInvokeEveryAnySql)
	}

	client := r.Group("api/client")
	client.Use(mid.AuthAll())
	{
		//客户端注册的相关接口
		cl := registered.Register{}
		//注册客户端
		client.POST("register", cl.Register)
		//查看客户端NPS配置
		client.GET("cl_npf_conf", cl.NpsConf)
		//查看客户端FRP配置
		client.GET("cl_frp_conf", cl.FrpConf)
	}
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(util.Port),
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	//signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server Exiting ...")
}
