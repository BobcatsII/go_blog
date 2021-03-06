package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"go_blog/models"
)

// 把每一步数据库操作都进行封装，
// 等待logic层根据业务需求调用。

const secret = "wwxad13wq5rk35r34113ef23f232f"

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对password的值进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("User Exist!")
	}
	return
}

// encryptPassword 加密密码
func encryptPassword(oPassword string) (string) {
	h := md5.New()
	h.Write([]byte(secret))
	// 下面将h.Sum生成的字符转换成16进制的字符串
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get (user, sqlStr, user.Username)
	if err == sql.ErrNoRows{
		return errors.New("用户不存在")
	}
	if err != nil{
		//查询数据库失败
		return err
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return errors.New("密码错误")
	}
	return
}