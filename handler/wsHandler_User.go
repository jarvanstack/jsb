package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WsUser struct {
	//发送信息的 chan
	ChanSendMessage chan string
	//websocket的连接
	conn *websocket.Conn
	room *WsRoom
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
		this.conn.WriteMessage(1, []byte(msg))
	}
}
//返回一个user对象
func NewWsUser(conn *websocket.Conn)*WsUser  {
	return &WsUser{
		conn: conn,
		ChanSendMessage: make(chan string),
		room: nil,
	}
}
//一致监听客户端的消息
func (this *WsUser)OnRead() {
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
		//未知消息
		this.unknownReadMsg(msg)
	}
}
//读取到未知消息
func (this *WsUser)unknownReadMsg(msg string)  {
	this.Write("我不是很懂您的意思哦("+msg+")")
}