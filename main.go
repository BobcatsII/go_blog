package main

import (
	"context"
	"fmt"
	"go_blog/controllers"
	"go_blog/dao/mysql"
	"go_blog/dao/redis"
	"go_blog/logger"
	"go_blog/pkg/snowflake"
	"go_blog/routes"
	"go_blog/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// go web开发 通用脚手架模板
func main() {
	// 1.加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v", err)
		return
	}

	// 2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v", err)
		return
	}
	defer zap.L().Sync() //把缓冲区的日志追加到日志文件中
	zap.L().Debug("logger init success...")

	// 3.初始化MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v", err)
		return
	}
	defer mysql.Close()

	// 4.初始化Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v", err)
		return
	}
	defer redis.Close()

	// x.雪花算法初始化
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v", err)
		return
	}
	// t.初始化gin框架内置的校验翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 5.注册路由
	r := routes.SetupRouter(settings.Conf.Mode)
	// 6.启动服务（含优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)                      // 创建一个接收信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
