package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	cLogin2 "github.com/uouuou/ServerManagerSystem/client/cLogin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"github.com/uouuou/ServerManagerSystem/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// ClientRouter 客户端路由
func ClientRouter() {
	gin.ForceConsoleColor()
	r := gin.New()
	r.NoRoute(mid.HandleNotFound)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(mid.Cors())
	v1 := r.Group("api/v1")
	v1.Use(mid.AuthAll())
	{
		v1.GET("token", cLogin2.GetToken)
	}
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(util.Port),
		Handler: r,
	}
	notAuth := r.Group("api/open")
	{
		//系统登录接口
		notAuth.POST("login", cLogin2.GetToken)

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
