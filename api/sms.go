package api

import (
	"fmt"
	"gin/service"
	"gin/tool"
	"github.com/gin-gonic/gin"
)

//注册登录发送短信
func sendSms(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	if len(mobile) != 11 {
		tool.ReturnFailure(ctx, 412, "发送验证码失败")
		return
	}
	//判断该电话号码是否已经被注册过
	flag, err := service.IsRegister(mobile)
	if err != nil {
		fmt.Println("查询电话号码出现错误", err)
		tool.ReturnFailure(ctx, 500, "发送验证码失败")
		return
	}
	//flag为false表示未注册
	if !flag {
		err = service.Register(mobile)
		if err != nil {
			fmt.Println("注册失败", err)
			tool.ReturnFailure(ctx, 500, "发送验证码失败")
			return
		}
	}
	err = service.SendSms(mobile)
	if err != nil {
		fmt.Println("验证码生成失败或接口调用失败", err)
		tool.ReturnFailure(ctx, 500, "发送验证码失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, "验证码发送成功")
}

//找回密码发送短信
func sendSmsForPwd(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	if len(mobile) != 11 {
		tool.ReturnFailure(ctx, 412, "电话号码格式错误")
		return
	}
	//判断该电话号码是否已经被注册过
	flag, err := service.IsRegister(mobile)
	if err != nil {
		fmt.Println("查询电话号码出现错误", err)
		tool.ReturnFailure(ctx, 500, "数据库查询错误")
		return
	}
	if !flag {
		tool.ReturnFailure(ctx, 412, "电话号码未注册")
		return
	}
	err = service.SendSms(mobile)
	if err != nil {
		fmt.Println("验证码生成失败或接口调用失败", err)
		tool.ReturnFailure(ctx, 500, "发送验证码失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, "验证码发送成功")
}
