package middleware

import "github.com/gin-gonic/gin"
//请求头和响应头的处理
func HeaderMiddlerware(ctx *gin.Context)  {

	ctx.Header("content-type","application/json")
	//放行
	ctx.Next()
}
