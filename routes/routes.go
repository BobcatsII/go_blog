package routes

import (
	"github.com/gin-gonic/gin"
	"go_blog/controllers"
	"go_blog/logger"
	"go_blog/middlewares"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式，不输出信息
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	//注册业务路由
	r.POST("/signup", controllers.SignUpHandler)
	r.POST("/login", controllers.LoginHandler)
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}


