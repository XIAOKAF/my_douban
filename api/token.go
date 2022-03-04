package api

import (
	"fmt"
	"gin/service"
	"gin/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func getRefreshToken(ctx *gin.Context) {
	refreshToken := ctx.PostForm("refreshToken")
	claims, err := service.ParseToken(refreshToken)
	//token类型错误
	if claims.Variety == "token" {
		tool.ReturnFailure(ctx, 500, "refreshToken种类错误")
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

	//无痛刷新
	//检验refreshToken的有效性
	//是否过期
	if time.Now().Before(claims.ExpireTime) {
		tool.ReturnFailure(ctx, 200, "RefreshToken is expired, please reLogin.")
		return
	}
	//认证refreshToken中的有效信息
	err, flag := service.CheckRefreshToken(claims.Mobile)
	if err != nil {
		fmt.Println("验证refreshToken失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	if flag == false {
		tool.ReturnFailure(ctx, 200, "无效refreshToken")
		return
	}
	//生成新的token,有效时间2分钟
	token, err := service.CreatToken(claims.Mobile, "token", 2)
	if err != nil {
		fmt.Println("token生成失败", err)
		tool.ReturnFailure(ctx, 500, "系统错误")
		return
	}
	//生成新的refreshToken，有效时间24小时
	refreshToken, err = service.CreatToken(claims.Mobile, "refreshToken", 1440)
	if err != nil {
		fmt.Println("refreshToken生成失败", err)
		tool.ReturnFailure(ctx, 500, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":         200,
		"token":        token,
		"refreshToken": refreshToken,
	})

}
