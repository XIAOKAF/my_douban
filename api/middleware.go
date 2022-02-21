package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors(ctx *gin.Context) {
	method := ctx.Request.Method
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("origin"))
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS") //设置服务器支持的所有跨域请求的方法
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	ctx.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
