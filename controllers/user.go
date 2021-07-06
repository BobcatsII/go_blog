package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go_blog/logic"
	"go_blog/models"
	"net/http"
)

func SignUpHandler(c *gin.Context) {
	// 1.获取参数 -> 参数校验
	p := new(models.ParamSignUp) //在models里提前定义好了参数结构体。
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误,直接返回响应
		zap.L().Error("SignUp with invalid param.", zap.Error(err))
		//判断errs是不是 validator.ValidationErrors 的可翻译类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(), //也可以定义一些业务码。
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)), //翻译一下错误，这个trans就是main函数里initTrans这一步初始化的翻译器
		})
		return
	}
	fmt.Println(p)
	// 2.业务处理（处理在logic层）
	if err := logic.SignUp(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	// 3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "注册成功.",
	})
}

func LoginHandler(c *gin.Context) {
	// 1.获取请求参数以及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断errs是不是 validator.ValidationErrors 的可翻译类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(), //也可以定义一些业务码。
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)), //翻译一下错误，这个trans就是main函数里initTrans这一步初始化的翻译器
		})
		return
	}
	// 2.业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}
	// 3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "登陆成功",
	})
}
