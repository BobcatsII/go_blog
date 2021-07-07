package jwt

import (
	"errors"
	"time"
)
import "github.com/dgrijalva/jwt-go"

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("ewk84lsdjfoi8jwnekl39nsdfjwe9orj1oands4fkljew7")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的 jwt.StandardClaims 只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体里面来
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个自己的声明的数据
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), //过期时间
			Issuer:    "go_blog",                                  //签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获取完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析Token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid{ //校验token
		return mc, nil
	}
	return nil, errors.New("invalid token.")
}