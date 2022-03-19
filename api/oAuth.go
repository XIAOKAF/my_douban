package api

import (
	"fmt"
	"gin/model"
	"gin/service"
	"gin/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

func oAuth(ctx *gin.Context) {
	code := ctx.PostForm("code")
	oauth := model.OAuth{
		ClientId:     "516**",
		ClientSecret: "c81**",
		RedirectUrl:  "http**",
	}
	t := model.Token{}
	oAuthUrl := service.GetUrl(oauth, code)
	token, err := service.GetToken(oAuthUrl, t)
	fmt.Println(err)
	if err != nil {
		fmt.Println("授权错误", err)
		tool.ReturnFailure(ctx, 500, "授权注册或登录失败")
		return
	}
	user := model.User{}
	user, err = service.GetUserInfo(token, user)
	if err != nil {
		fmt.Println("解析用户信息错误", err)
		tool.ReturnFailure(ctx, 500, "授权注册或登录失败")
		return
	}
	//查询该用户是否注册过
	err, flag := service.CheckRegisterByUsername(user)
	if err != nil {
		fmt.Println("查询用户是否注册错误", err)
		tool.ReturnFailure(ctx, 500, "授权注册或登录失败")
		return
	}
	//未查询到该用户
	if !flag {
		err = service.RegisterByOAuth(user)
		if err != nil {
			fmt.Println("注册错误", err)
			tool.ReturnFailure(ctx, 500, "授权注册或登录失败")
			return
		}
	}
	//记住登录状态
	//生成token,有效时间2分钟
	LoginToken, err := service.CreatToken(user.Username, "token", 30)
	if err != nil {
		fmt.Println("token生成失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//生成refreshToken，有效时间24小时
	refreshToken, err := service.CreatToken(user.Username, "refreshToken", 1440)
	if err != nil {
		fmt.Println("refreshToken生成失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":         200,
		"token":        LoginToken,
		"refreshToken": refreshToken,
	})
}
