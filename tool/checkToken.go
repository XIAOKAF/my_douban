package tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CheckToken(ctx *gin.Context, err error) bool {
	if err != nil {
		if err.Error()[:16] == "token is expired" {
			ReturnFailure(ctx, 200, "token is expired")
			return false
		}
		fmt.Println("token解析错误", err)
		ReturnFailure(ctx, 500, "token解析错误")
		return false
	}
	return true
}
