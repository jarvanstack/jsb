package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"jsb/util/my_token"
	"log"
	"net/http"
	"sync"
	"time"
)


type WsServer struct {
	//websocket 的房间
	WsRooms map[int]*WsRoom
	//websocket User 池
	WsUsers map[int64]*WsUser
	//读写锁
	LockWsUsers sync.RWMutex
}

func NewWsServer()*WsServer  {
	return &WsServer{
		WsRooms: make(map[int]*WsRoom),
		WsUsers: make(map[int64]*WsUser),
	}
}
//进来一个
func (this *WsServer)WsHandler( ctx *gin.Context) {
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
	fmt.Printf("新的websocket=%#v\n", conn.RemoteAddr().String())
	//1.创建一个新的 User
	token := ctx.GetHeader("token")
	tokenUser, err := my_token.GetUser(token)
	wsUser := NewWsUser(conn)
	//放入池子中
	this.LockWsUsers.Lock()
	this.WsUsers[tokenUser.UserId] = wsUser
	this.LockWsUsers.Unlock()
	//开始监听写入
	go wsUser.OnWrite()
	//监听读取客户端信息
	go wsUser.OnRead()

}


