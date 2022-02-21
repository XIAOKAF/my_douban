package api

import (
	"fmt"
	"gin/model"
	"gin/service"
	"gin/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
)

//短信注册或登录
func registerOrLoginByVerifyCode(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	verifyCode := ctx.PostForm("verifyCode")
	postTime := time.Now()
	if len(mobile) != 11 {
		tool.ReturnFailure(ctx, 412, "电话号码格式错误")
		return
	}
	//判断该电话号码是否已经被注册过
	flag, err := service.IsRegister(mobile)
	if err != nil {
		fmt.Println("查询电话号码出现错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//电话号码未被注册
	if !flag {
		tool.ReturnFailure(ctx, 412, "电话号码错误")
		return
	}
	//验证码是否正确且未过期
	code, sendTime, err := service.SelectVerifyCodeAndSendTime(mobile)
	if err != nil {
		fmt.Println("查询验证码和时间错误", err)
		tool.ReturnFailure(ctx, 500, "登录失败")
		return
	}
	mm, err := time.ParseDuration("30m")
	if err != nil {
		fmt.Println("时间差定义错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	latestTime := sendTime.Add(mm)
	if latestTime.Before(postTime) {
		tool.ReturnFailure(ctx, 408, "验证码已过期")
		return
	}
	if verifyCode != code {
		tool.ReturnFailure(ctx, 410, "验证码错误")
		return
	}
	//生成token,有效时间2分钟
	token, err := service.CreatToken(mobile, "token", 30)
	if err != nil {
		fmt.Println("token生成失败", err)
		tool.ReturnFailure(ctx, 500, "系统错误")
		return
	}
	//生成refreshToken，有效时间24小时
	refreshToken, err := service.CreatToken(mobile, "refreshToken", 1440)
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

func improvePersonalInfo(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	username := ctx.PostForm("username")
	password := ctx.PostForm("pwd")
	selfIntroduction := ctx.PostForm("selfIntroduction")
	if token == "" {
		tool.ReturnFailure(ctx, 200, "No Token")
		return
	}
	//token过期
	claims, err := service.ParseToken(token)
	flag := tool.CheckToken(ctx, err)
	if !flag {
		return
	}
	//token类型错误
	if claims.Variety == "refreshToken" {
		tool.ReturnFailure(ctx, 200, "token类型错误")
		return
	}
	flag = tool.CheckToken(ctx, err)
	if !flag {
		return
	}

	//用户名已被使用
	err, flag = service.SelectUsernameByMobile(claims.Mobile)
	if err != nil {
		fmt.Println("查询用户名错误", err)
		tool.ReturnFailure(ctx, 500, "数据库查询错误")
		return
	}
	if flag {
		tool.ReturnFailure(ctx, 200, "用户名已被使用")
		return
	}
	//密码长度不符合要求
	if len(password) > 15 || len(password) < 6 {
		tool.ReturnFailure(ctx, 200, "密码的长度必须在6到15位之间")
		return
	}
	//密码过于简单
	pwd, err := regexp.Compile(`[A-Za-z][0-9]`)
	if err != nil {
		fmt.Println("密码解析错误", err)
		tool.ReturnFailure(ctx, 200, "密码设置错误")
		return
	}
	p := string(pwd.Find([]byte(password)))
	if len(p) < 2 {
		tool.ReturnFailure(ctx, 200, "密码必须包含字母与数字")
		return
	}
	//自我介绍过长
	if len(selfIntroduction) > 250 {
		tool.ReturnFailure(ctx, 200, "自我介绍过长")
		return
	}
	user := model.User{
		Mobile:           claims.Mobile,
		Username:         username,
		Password:         password,
		SelfIntroduction: selfIntroduction,
	}
	err = service.ImprovePersonalInfo(user)
	if err != nil {
		fmt.Println("更新个人信息失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	tool.ReturnSuccess(ctx, 200, "更新成功")
}

func getPersonalInfo(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	err, user := service.SelectInfoByMobile(mobile)
	if err != nil {
		fmt.Println("获取个人信息失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	tool.ReturnSuccess(ctx, 200, user)
}

func loginByPassword(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("pwd")
	user := model.User{
		Mobile:   mobile,
		Password: password,
	}
	err, pwd := service.SelectPasswordByMobile(user)
	if err != nil {
		fmt.Println("查询密码错误", err)
		tool.ReturnFailure(ctx, 500, "登录失败")
		return
	}
	if pwd != user.Password {
		tool.ReturnFailure(ctx, 401, "密码错误")
		return
	}
	//生成token,有效时间2分钟
	token, err := service.CreatToken(mobile, "token", 2)
	if err != nil {
		fmt.Println("token生成失败", err)
		tool.ReturnFailure(ctx, 500, "系统错误")
		return
	}
	//生成refreshToken，有效时间24小时
	refreshToken, err := service.CreatToken(mobile, "refreshToken", 1440)
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

func changePasswordByVerifyCode(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	newPwd := ctx.PostForm("newPwd")
	verifyCode := ctx.PostForm("verifyCode")
	postTime := time.Now()
	user := model.User{
		Mobile:     mobile,
		Password:   newPwd,
		VerifyCode: verifyCode,
	}
	code, sendTime, err := service.SelectVerifyCodeAndSendTime(mobile)
	if err != nil {
		fmt.Println("查询验证码和时间错误", err)
		tool.ReturnFailure(ctx, 500, "登录失败")
		return
	}
	mm, err := time.ParseDuration("30m")
	if err != nil {
		fmt.Println("时间差定义错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	latestTime := sendTime.Add(mm)
	if latestTime.Before(postTime) {
		tool.ReturnFailure(ctx, 401, "验证码已过期")
		return
	}
	if verifyCode != code {
		tool.ReturnFailure(ctx, 200, "验证码错误")
		return
	}
	//密码长度不符合要求
	if len(newPwd) > 15 || len(newPwd) < 6 {
		tool.ReturnFailure(ctx, 200, "密码的长度必须在6到15位之间")
		return
	}
	//密码过于简单
	pwd, err := regexp.Compile(`[A-Za-z][0-9]`)
	if err != nil {
		fmt.Println("密码解析错误", err)
		tool.ReturnFailure(ctx, 500, "密码解析错误")
		return
	}
	n := string(pwd.Find([]byte(newPwd)))
	if len(n) < 2 {
		tool.ReturnFailure(ctx, 200, "密码必须同时包含数字和英文字母")
		return
	}
	err, oldPwd := service.SelectPasswordByMobile(user)
	if err != nil {
		fmt.Println("数据库查询错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	if oldPwd == newPwd {
		tool.ReturnFailure(ctx, 200, "密码不能和原密码相同")
		return
	}
	//更新密码
	err = service.ChangePassword(user)
	if err != nil {
		fmt.Println("密码重置失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	tool.ReturnSuccess(ctx, 200, "密码重置成功")
}
