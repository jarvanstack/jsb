package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

//
func WsHandler( ctx *gin.Context) {
	//w http.ResponseWriter, r *http.Request
	w := ctx.Writer
	r := ctx.Request
	//创建 upgrader
	upgrade := &websocket.Upgrader{
		//超时时间一个小时
		HandshakeTimeout: time.Minute * 60,
		//读写缓存池
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		//请求检查函数，用于统一的链接检查，以防止跨站点请求伪造。如果不检查，就设置一个返回值为true的函数。
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//升级连接 并且拿到 websocket conn
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("string=%s\n", "升级连接err")
	}
	fmt.Printf("conn.RemoteAddr()=%#v\n", conn.RemoteAddr().String())
	//使用 conn 读取和写入内容
	go ReadFromWS(conn)
	go WriteToWS(conn)
}

//读取消息并打印
func ReadFromWS(conn *websocket.Conn) {
	for {
		message, bytes, err := conn.ReadMessage()
		//断线不自动重新连接
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("type=%#v\n", message)
		fmt.Printf("message=%s\n", string(bytes))

	}
}
func WriteToWS(conn *websocket.Conn) {
	for  {
		select {
		case <-time.After(time.Second * 10):
			conn.WriteMessage(1,[]byte("服务端发送信息"))
		}

	}
}
