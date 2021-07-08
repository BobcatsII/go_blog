package middlewares

import (
	"github.com/gin-gonic/gin"
	"go_blog/controllers"
	"go_blog/pkg/jwt"
	"strings"
)

//此文件为认证文件
// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的 Authorization 中，并使用Bearer开头
		// Authorization: Bearer xxxxxx.xxxxx.xxxx 或 X-TOKEN: xxx.xxx.xxx 格式
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割Authorization, 去掉Bearer, 取得Token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controllers.CtxtUserIDKey, mc.UserID)
		c.Next() // 后续的处理请求的函数追踪可以用过c.Get("username")来获取当前请求的用户信息
	}
}
