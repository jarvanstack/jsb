package my_gorm

import (
	"fmt"
	"jsb/model/entity"
	"jsb/util/snowflake"
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
func TestCreateSysResult(t *testing.T) {
	err := DB.AutoMigrate(&entity.SysResult{})
	if err != nil {
		fmt.Printf("err=%#v\n", err)
	}
}
func TestInsertResult(t *testing.T) {
	sql := `
	INSERT INTO jsb.sys_results 
	( id, a_user_id, b_user_id, a_result, b_result, round_num )
	VALUES
	(?,?1,'?','?','?',?)
	`
	id := snowflake.NextId()

	s := DB.Exec(sql, id, 1, 1, 1, 1, 1).Error
	fmt.Printf("s=%#v\n", s)

}