package handler

import (
	"jsb/util/snowflake"
	"log"
	 "jsb/util/my_gorm"
)

type WsRoom struct {
	//第几局
	RoundNum int
	//房间号
	RoomId int
	//userId 用户定值
	AWsUser *WsUser
	//手势 j、s、b
	AResult string
	//userId 用户定值
	BwsUser *WsUser
	//手势 j、s、b
	BResult string
}

func NewWsRoom(roomId int)*WsRoom  {
	return &WsRoom{
		AWsUser: nil,
		BwsUser: nil,
		AResult: "",
		BResult: "",
		RoomId: roomId,
		RoundNum: 0,
		}
}
//加入user到Room里面
func (this *WsRoom) AddUser(wsUser *WsUser) bool  {
	if this.AWsUser == nil {
		this.AWsUser = wsUser
	}else if this.BwsUser == nil {
		this.BwsUser = wsUser
	}else {
		log.Printf("string=%#v\n", "该房间用户已经满了")
		return false
	}
	return true
}
//添加Result
func (this *WsRoom) AddResult(wsUser *WsUser,gesture string) bool  {
	if this.AWsUser == wsUser {
		this.AResult = gesture
	}else if this.BwsUser == wsUser {
		this.BResult = gesture
	}else {
		log.Printf("string=%#v\n", "该房间结果已经满了")
		return false
	}
	return true
}
//user数量
func (this *WsRoom)NumOfWsUser()int  {
	numOfWsUser := 0;
	if nil != this.AWsUser  {
		numOfWsUser++
	}
	if nil != this.BwsUser  {
		numOfWsUser++
	}
	return numOfWsUser
}
//数量Result
func (this *WsRoom)NumOfWsResult()int  {
	numOfWsResult := 0;
	if this.AResult != "" {
		numOfWsResult++
	}
	if this.BResult != "" {
		numOfWsResult++
	}
	return numOfWsResult
}
func (this *WsRoom)saveToMysqlAndCleanResult()  {
	log.Printf("string=%#v\n", "储存结果并再来一局")
	//save to mysql
	a_user_id := this.AWsUser.UserId
	b_user_id := this.BwsUser.UserId
	a_result := this.AResult
	b_result := this.BResult
	round_num := this.RoundNum
	id := snowflake.NextId()

	sql := `
	INSERT INTO jsb.sys_results 
	( id, a_user_id, b_user_id, a_result, b_result, round_num )
	VALUES
	(?,?,?,?,?,?)
`	//clean result
	this.AResult = ""
	this.BResult = ""
	err := my_gorm.DB.Exec(sql, id, a_user_id, b_user_id, a_result, b_result, round_num).Error
	if err != nil {
		log.Printf("string=%#v\n", "插入失败")
	}
}