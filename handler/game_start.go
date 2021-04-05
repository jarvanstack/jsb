package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func GameStartHandler( ctx *gin.Context)  {
	fmt.Printf("string=%s\n", "进入方法")
	time.Sleep(time.Second*20)
	ctx.Writer.Write([]byte("阻塞结束"))
}
