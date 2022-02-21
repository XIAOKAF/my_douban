package api

import (
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()

	engine.Use(Cors) //处理跨域问题

	userGroup := engine.Group("/user")
	{
		userGroup.POST("/registerOrLoginByVerifyCode", registerOrLoginByVerifyCode) //验证码注册或登录
		userGroup.POST("/loginByPassword", loginByPassword)                         //密码登录
		userGroup.POST("/improvePersonalInfo", improvePersonalInfo)                 //完善个人信息
		userGroup.GET("/getPersonalInfo", getPersonalInfo)                          //查看个人信息
		userGroup.POST("/changePasswordByVerifyCode", changePasswordByVerifyCode)   //修改密码
	}

	homeGroup := engine.Group("/homepage")
	{
		homeGroup.GET("/hotShowing", hotShowing)                         //正在热映
		homeGroup.GET("/recentHotMovie", recentHotMovie)                 //最近热门电影
		homeGroup.GET("/recentHotTeleplay", recentHotTeleplay)           //一周口碑榜
		homeGroup.GET("/weeklyPraise", weeklyPraise)                     //最近热门电视剧
		homeGroup.GET("/hotRecommendation", hotRecommendation)           //热门推荐
		homeGroup.POST("/selectMoviesByKeyWords", SelectMovieByKeyWords) //关键字搜索电影
	}

	movieGroup := engine.Group("/movieDetails")
	{
		movieGroup.GET("/selectMovieDetails", selectMovieDetails) //影片详情
		movieGroup.POST("/postShortComment", PostShortComment)    //发布短评
		movieGroup.GET("/getComment", SelectComment)              //获取短评
	}

	engine.GET("/selectCelebrityDetails", getCelebrityDetails) //影人详情

	smsGroup := engine.Group("/sms")
	{
		smsGroup.POST("/sendSms", sendSms)             //登陆注册的短信接口
		smsGroup.POST("/sendSmsForPwd", sendSmsForPwd) //找回密码时的短信接口
	}

	engine.POST("/refreshToken", getRefreshToken) //通过refreshToken生成新的token

	engine.Run(":8090")
}
