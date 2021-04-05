package my_gorm

import (
	"fmt"
	"testing"
)
type Result struct {
	ID   int
	Name string
	Age  int
}
type SysUser struct {
	UserId string `json:"user_id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar string `json:"avatar"`
}

func createTable()  {
	err := DB.AutoMigrate(&Result{})
	if err != nil {
		fmt.Printf("err=%#v\n", err)
	}
}
func Test原始SQL(t *testing.T) {
	var user SysUser
	sql := `
	SELECT
		sys_user.user_id, 
		sys_user.username, 
		sys_user.email, 
		sys_user.password, 
		sys_user.nickname, 
		sys_user.avatar
	FROM
		sys_user
	WHERE
		sys_user.password = ? AND
		sys_user.username = ?
	limit 1
	`
	DB.Raw(sql,"admin","admin").Scan(&user)
	fmt.Printf("user=%#v\n", user)
}
