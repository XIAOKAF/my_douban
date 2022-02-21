package api

import (
	"fmt"
	"gin/service"
	"gin/tool"
	"github.com/gin-gonic/gin"
)

func getRefreshToken(ctx *gin.Context) {
	refreshToken := ctx.PostForm("refreshToken")
	claims, err := service.ParseToken(refreshToken)
	//token类型错误
	if claims.Variety == "token" {
		tool.ReturnFailure(ctx, 500, "refreshToken错误")
		return
	}
	if err != nil {
		if err.Error() == "token is expired" {
			tool.ReturnFailure(ctx, 401, "refreshToken已过期，请重新登录")
			return
		}
		fmt.Println("获取refreshToken失败", err)
		tool.ReturnFailure(ctx, 500, "系统错误")
	}
	/*
		//无痛刷新
			//生成token,有效时间2分钟
			token, err := service.CreatToken(claims.mobile,"token",2)
			if err != nil {
				fmt.Println("token生成失败", err)
				tool.ReturnFailure(ctx, "系统错误")
				return
			}
			//生成refreshToken，有效时间24小时
			refreshToken, err := service.CreatToken(claims.mobile,"refreshToken",1440)
			if err != nil {
				fmt.Println("refreshToken生成失败", err)
				tool.ReturnFailure(ctx, "系统错误")
				return
			}
			ctx.JSON(http.StatusOK,gin.H{
				"code": 200,
				"token": token,
				"refreshToken": refreshToken,
			})
	*/
}
