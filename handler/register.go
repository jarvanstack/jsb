package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jsb/model/ato"
	"jsb/util/my_gorm"
	"jsb/util/my_restful"
	"jsb/util/snowflake"
	"log"
)




func RegisterHandler( ctx *gin.Context)  {
	//拿到账号和密码
	var login ato.Login
	err2 := ctx.ShouldBindJSON(&login)
	if err2 != nil {
		log.Printf("string=%#v\n", "register json 绑定错误")
	}
	fmt.Printf("login=%#v\n", login)
	//在MySQL验证
	sql := "INSERT INTO `jsb`.`sys_user`(`user_id`, `username`, `password`) VALUES (?, ?, ?)"
	db := my_gorm.DB
	err := db.Exec(sql, snowflake.NextId(), login.Username, login.Password).Error
	//如果登录失败
	if err != nil {
		ctx.Writer.Write(my_restful.Fail("该用户已经注册了"))
		return
	}
	ctx.Writer.Write(my_restful.Ok("注册成功"))
}
