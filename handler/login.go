package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jsb/model/ato"
	"jsb/util/my_gorm"
	"jsb/util/my_restful"
	"jsb/util/my_token"
	"jsb/model/entity"
)




func LoginHandler( ctx *gin.Context)  {
	//拿到账号和密码
	var login ato.Login
	ctx.ShouldBindJSON(&login)
	fmt.Printf("login=%#v\n", login)
	//在MySQL验证
	var user entity.SysUser
	sql := `
	SELECT
		sys_user.user_id,
		sys_user.username, 
		sys_user.email,
		sys_user.nickname,
		sys_user.avatar
	FROM
		sys_user
	WHERE
		sys_user.password = ? AND
		sys_user.username = ?
	limit 1
	`
	db := my_gorm.DB
	db.Raw(sql,login.Password,login.Username).Scan(&user)
	fmt.Printf("user=%#v\n", user)
	//如果登录失败
	if user.UserId == 0 {
		ctx.Writer.Write(my_restful.Fail("登录失败"))
		return
	}
	token, err := my_token.GetToken(user)
	if err != nil {
		fmt.Printf("string=%s\n", "token err")
	}
	ctx.Header("token",token)
	ctx.Writer.Write(my_restful.Ok(user))
}
