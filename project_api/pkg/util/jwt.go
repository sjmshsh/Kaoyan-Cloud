package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtSecret = []byte("I love zyq!") // 声明签名信息

// Claims 自定义有效载荷
type Claims struct {
	Uid                uint   `json:"uid"`
	Username           string `json:"username"`
	jwt.StandardClaims        // StandardClaims结构体实现了Claims接口(Valid()函数)
}

// GenerateToken 签发token（调用jwt-go库生成token）
func GenerateToken(uid uint, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * 24)
	claims := Claims{
		Uid:      uid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: nowTime.Unix(),    // 签名生效时间
			ExpiresAt: expireTime.Unix(), // 签名过期时间
			Issuer:    "Lxy",             // 签名颁发者
		},
	}
	// 指定编码算法为jwt.SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 返回一个token结构体指针(*Token)
	//tokenString, err := token.SigningString(JwtSecret)
	//return tokenString, err
	return token.SignedString(JwtSecret)
}

// ParseToken token解码
func ParseToken(tokenString string) (*Claims, error) {
	// 输入用户token字符串,自定义的Claims结构体对象,以及自定义函数来解析token字符串为jwt的Token结构体指针
	//Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	//func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 将token中的claims信息解析出来,并断言成用户自定义的有效载荷结构
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token不可用")
}
