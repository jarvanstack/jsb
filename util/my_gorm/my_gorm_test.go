package my_gorm

import (
	"bytes"
	"fmt"
	"jsb/model/entity"
	"jsb/util/snowflake"
	"strconv"
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

type TestResult struct {
	AUserId int64
	BUserId int64
	AResult string
	BResult string
	RoundNum int
}
func TestGetResult(t *testing.T) {
	userId := int64(132414)
	sql := `
	SELECT
		a_user_id, 
		b_user_id, 
		a_result, 
		b_result, 
		round_num
	FROM
		sys_results
	WHERE
		sys_results.a_user_id = ? OR
		sys_results.b_user_id = ?
	limit 10
`
	var results []TestResult
	err := DB.Raw(sql, userId, userId).Scan(&results)
	if err != nil {
		fmt.Printf("err=%#v\n", err)
	}
	var bufferBytes bytes.Buffer
	for index,result := range results {
		var opponentUserId int64
		if userId==result.AUserId {
			opponentUserId = result.BUserId
		}else {
			opponentUserId = result.AUserId
		}
		bufferBytes.WriteString(strconv.Itoa(index + 1))
		bufferBytes.WriteString(". vs ")
		bufferBytes.WriteString(strconv.FormatInt(opponentUserId, 10))
		//拿到自己的结果和对方的结果
		var myResult,opponentResult string
		if userId==result.AUserId{
			myResult = result.AResult
			opponentResult = result.BResult
		}else {
			myResult = result.BResult
			opponentResult = result.AResult
		}
		bufferBytes.WriteString(" 你:"+myResult)
		bufferBytes.WriteString(" 对方:"+opponentResult)
		bufferBytes.WriteString(";")
	}
	fmt.Printf("拼接结果=%#v\n", bufferBytes.String())




}