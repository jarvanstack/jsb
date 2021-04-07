package main

import (
	"github.com/gin-gonic/gin"
	"jsb/handler"
	"jsb/middleware"
)


func main() {
	engine := gin.Default()
	//中间件
	engine.Use(middleware.HeaderMiddlerware)
	//登录
	engine.POST("/jsb/login",handler.LoginHandler)
	engine.POST("/jsb/register",handler.RegisterHandler)
	//构建Websocket server
	wsServer := handler.NewWsServer()
	engine.GET("/jsb/send-message",wsServer.WsHandler)
	engine.Run("localhost:8080")
}

