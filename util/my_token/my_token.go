package my_token

import (
	"jsb/model/entity"
	"time"
	"github.com/dgrijalva/jwt-go"
)

//私钥
var privateKeyByte []byte
//过期时间
var  expire time.Duration
//签发者
var issuer string
//签名主题
var subject string

//初始化默认值
func init()  {
	privateKeyByte = []byte("这是一个私钥")
	expire =  7* 24 * time.Hour
	issuer = "ftcloud"
	subject = "user my_token"

}

type ClaimsObj struct {
	//用户id
	User entity.SysUser
	jwt.StandardClaims
}

func GetToken(user  entity.SysUser)(string,error)  {
	expireTime := time.Now().Add(expire)
	claims := &ClaimsObj{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:   issuer,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(privateKeyByte)
}
func GetUser(tokenString string)(user  entity.SysUser,err error)  {
	claimsObj, err := ParseToken(tokenString)
	return claimsObj.User,err
}
func ParseToken(tokenString string) ( *ClaimsObj, error) {
	claimsObj := &ClaimsObj{}
	_, err := jwt.ParseWithClaims(tokenString, claimsObj, func(token *jwt.Token) (i interface{}, err error) {
		return privateKeyByte, nil
	})
	return claimsObj, err
}