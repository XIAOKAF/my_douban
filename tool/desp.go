package tool

import "github.com/gin-gonic/gin"

func ReturnFailure(ctx *gin.Context, code int, info interface{}) {
	ctx.JSON(code, gin.H{
		"code": code,
		"info": info,
	})
}

func ReturnSuccess(ctx *gin.Context, code int, info interface{}) {
	ctx.JSON(code, gin.H{
		"code": code,
		"info": info,
	})
}
