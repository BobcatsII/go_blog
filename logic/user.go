package logic

import (
	"go_blog/dao/mysql"
	"go_blog/models"
	"go_blog/pkg/snowflake"
)

// logic 层是存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 0.判断用户是否存在于数据库，如果存在表示已注册，不继续往下走；
	// 如果不存在则往下走
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 这一步是数据库查询出错
		return err
	}
	// 1.生成UID
	userID:=snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID: userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 2.设置密码，并保存到数据库(涉及到dao层)
	return mysql.InsertUser(user)

	// 3.redis.xxx ...
}

func Login(p *models.ParamLogin) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)

}