package handler

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"jsb/model/entity"
	"jsb/util/my_gorm"
	"jsb/util/my_string_util"
	"strconv"
	"sync"
)

type WsUser struct {
	UserId int64
	//发送信息的 chan
	ChanSendMessage chan string
	//解决并发写入websocket问题
	wsConnLock sync.RWMutex
	//websocket的连接
	conn   *websocket.Conn
	room   *WsRoom
	server *WsServer
}

//写入数据回去
func (this *WsUser) Write(msg string) {
	this.ChanSendMessage <- msg
}

//我们需要一致监听然后写入
func (this *WsUser) OnWrite() {
	for true {
		msg := <-this.ChanSendMessage
		if msg == "closeOnWrite" {
			return
		}
		this.wsConnLock.Lock()
		this.conn.WriteMessage(1, []byte(msg))
		this.wsConnLock.Unlock()
	}
}

//返回一个user对象
func (this *WsServer) newWsUser(conn *websocket.Conn, userId int64) *WsUser {
	return &WsUser{
		conn:            conn,
		ChanSendMessage: make(chan string),
		room:            nil,
		UserId:          userId,
		server:          this,
	}
}

//一致监听客户端的消息
func (this *WsUser) OnRead() {
	for {
		messageType, bytes, err := this.conn.ReadMessage()
		//断线不自动重新连接
		if err != nil {
			fmt.Printf("断开websocket=%s\n", this.conn.RemoteAddr().String())
			this.conn.Close()
			return
		}
		msg := string(bytes)
		fmt.Printf("messageType=%#v\n", messageType)
		fmt.Printf("读取到内容=%s\n", msg)

		//1.接受到 jsb 创建一个房间
		if msg == "jsb" {
			length := len(this.server.WsRooms)
			this.server.LockWsRooms.Lock()
			this.server.WsRooms[length] = NewWsRoom(length)
			this.server.LockWsRooms.Unlock()
			this.Write("您的房间号是:「" + strconv.Itoa(length) + "」请您和您的对手输入房间号进入房间比如「13」")
			continue
		}
		//2.进入房间
		if my_string_util.IsNum(msg) {
			//1.看看有没有这个房间
			roomNum, _ := strconv.Atoi(msg)
			room := this.server.WsRooms[roomNum]
			if room == nil {
				this.Write("没有这个房间「" + msg + "」请重新输入房间号或者输入「jsb」创建新的房间")
				continue
			}
			//这个房间找得到不？？
			numOfWsUser := this.server.WsRooms[roomNum].NumOfWsUser()
			//如果房间么没有满
			if  numOfWsUser==0{
				//2.如果有这个房间就加入房间，并返回加入成功
				//2.1加入user到room
				room.AddUser(this)
				//2.2加入room到user
				this.room = room
				this.Write("进入房间成功，请让您的队友尽快进入房间")
			}else if numOfWsUser==1 {
				room.AddUser(this)
				this.room = room
				this.Write("进入房间成功，您的对手也已经准备好，输入「j」「s」「b」开始对决")
				this.getOpponent().Write("您的对手也已经准备好，输入「j」「s」「b」开始对决")
			} else {
				this.Write("房间" + strconv.Itoa(roomNum) + "已经满了，请输入「jsb」开始一个新的房间，或者「j」、「s」、「b」代表剪刀、石头、布开始游戏对决")

			}

			continue
		}
		//3.输入剪刀石头布开始对决
		if Isjsb(msg) {
			numOfWsResult := this.room.NumOfWsResult()
			//3.1储存结果并通知自己出手成功，通知对方自己已经出手.
			if numOfWsResult==0 {
				//添加手势到结果
				this.room.AddResult(this,msg)
				//通知自己出手成功
				this.Write("出招成功「"+msg+"」")
				//通知对方
				opponent := this.getOpponent()
				opponent.Write("对方已经出招成功，请尽快出招")
			} else if numOfWsResult==1 {
				//3.2如果结果已经有1个，就添加第二个,并公布答案并协程转储存到MySQL.并协程清理房间的result
				//添加手势到结果
				this.room.AddResult(this,msg)
				//公布结果到对应的用户
				opponentResult := this.getOpponentResult()
				//局数 + 1
				this.room.RoundNum ++
				this.Write("你的出招「"+msg+"」对方的出招「"+opponentResult+"」--局数:"+strconv.Itoa(this.room.RoundNum))
				this.getOpponent().Write("你的出招「"+opponentResult+"」对方的出招「"+msg+"」--局数:"+strconv.Itoa(this.room.RoundNum))
				//协程转储到MySQL 并清理.
				go this.saveToMysqlAndCleanResult()
				this.Write("新的一局,「j」、「s」、「b」代表剪刀、石头、布开始继续对决")
				this.getOpponent().Write("新的一局,「j」、「s」、「b」代表剪刀、石头、布继续游戏对决")
			}else {
				//3.3 如果已经满了就返回错误
				this.Write("错误:结果已经存在")
			}
			continue
		}else if msg=="jl" {
			//展示记录
			userId := this.UserId
			results := getResultsById(userId)
			this.Write(results)
			continue
		}


		// x.未知消息
		this.unknownReadMsg(msg)
	}
}
//SQL查询Result记录.
func getResultsById(userId int64)string  {
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
	var results []entity.SysResult
	err := my_gorm.DB.Raw(sql, userId, userId).Scan(&results)
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
	return bufferBytes.String()
}

//读取到未知消息
func (this *WsUser) unknownReadMsg(msg string) {
	this.Write("我不是很懂您的意思哦(" + msg + ")")
}

//判断是否是 j、s、或者b
func Isjsb(str string) bool {
	switch str {
	case "J":
		return true
	case "S":
		return true
	case "B":
		return true
	case "j":
		return true
	case "s":
		return true
	case "b":
		return true
	default:
		return false
	}
}
//获取对手
func (this *WsUser)getOpponent()*WsUser  {
	opponent := this.room.AWsUser
	if this == opponent {
		opponent = this.room.BwsUser
	}
	return opponent
}
func (this *WsUser)getOpponentResult()string  {
	opponent := this.room.AWsUser
	opponentResult := this.room.AResult
	if this == opponent {
		opponentResult = this.room.BResult
	}
	return opponentResult
}
func (this *WsUser)saveToMysqlAndCleanResult()  {
	this.room.saveToMysqlAndCleanResult()
}