package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jsb/model/ato"
	"jsb/util/my_restful"
)

func GameStartHandler( ctx *gin.Context)  {
	//拿到message
	fmt.Printf("string=%s\n", "进入方法GameStartHandler")
	message := ato.Message{}
	ctx.ShouldBindJSON(&message)
	fmt.Printf("message=%#v\n", message)
	//处理游戏开始
	if message.Message=="jsb" {
		fmt.Printf("string=%s\n", "处理游戏开始")
		//开始游戏
		ctx.Writer.Write(my_restful.Ok("游戏开始你的房间号为"))
	}
	ctx.Writer.Write(my_restful.Ok("我不是很懂你的意思哦"))
}
