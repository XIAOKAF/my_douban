package api

import (
	"fmt"
	"gin/model"
	"gin/service"
	"gin/tool"
	"github.com/gin-gonic/gin"
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
	err = service.CheckRegisterByUsername(user)
	if err != nil {
		fmt.Println("查询用户是否注册错误", err)
		tool.ReturnFailure(ctx, 500, "授权注册或登录失败")
		return
	}
	//授权接口中注册

	//记住登录状态
}
