package handler

import (
	"fmt"
	"jsb/util/my_gorm"
	"jsb/util/snowflake"
	"testing"
)

func TestRegister(t *testing.T) {
	sql := "INSERT INTO `jsb`.`sys_user`(`user_id`, `username`, `password`) VALUES (?, ?, ?)"

	err := my_gorm.DB.Exec(sql, snowflake.NextId(), "admin3", "admin3").Error
	if err != nil {
		fmt.Printf("err=%#v\n", err)
	}
}

