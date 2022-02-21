package api

import (
	"fmt"
	"gin/model"
	"gin/service"
	"gin/tool"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
)

func getCelebrityDetails(ctx *gin.Context) {
	celebrityId := ctx.PostForm("celebrityId")
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/celebrity/"+celebrityId+"/", nil)
	if err != nil {
		fmt.Println("请求失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36 Edg/98.0.1108.43")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "ll=\"108309\"; bid=gtXWf_veA68; push_doumail_num=0; push_noty_num=0; __utmv=30149280.21338; __yadk_uid=ushdXmeu6hPz8VL7FLqtxfB7XyoIDI0K; __gads=ID=077fbb5baee62c93-227b63397dd000ae:T=1644204395:RT=1644204395:S=ALNI_MawllVEXeybHBZ7lVBDC1sUfvJgCg; ct=y; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; dbcl2=\"213387422:/wRfw1lR4Aw\"; Hm_lvt_16a14f3002af32bf3a75dfe352478639=1644762356; ck=O7li; __utmc=30149280; __utmc=223695111; __utmz=223695111.1645065903.38.19.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ps=y; __utmz=30149280.1645102423.42.20.utmcsr=cn.bing.com|utmccn=(referral)|utmcmd=referral|utmcct=/; _pk_ref.100001.4cf6=[\"\",\"\",1645197224,\"https://www.douban.com/\"]; ap_v=0,6.0; __utma=30149280.393236484.1643783658.1645102423.1645197229.43; __utma=223695111.2053381400.1643783658.1645065903.1645197229.39; _pk_id.100001.4cf6=63f53dbc9e9bd5ef.1643783657.39.1645197902.1645066135.")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}
	defer resp.Body.Close()

	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//查询该演员是否录入
	err, flag := service.SelectCelebrityById(celebrityId)
	if err != nil {
		fmt.Println("查询演员失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//该演员已经录入，则直接返回信息
	if flag {
		//获取演员信息
		err, celebrity := service.SelectCelebrityDetails(celebrityId)
		if err != nil {
			fmt.Println("查询演员信息错误", err)
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
		//获取演员照片
		err, celebrity.Photos = service.SelectPhotos(celebrityId)
		if err != nil {
			fmt.Println("查询演员照片失败", err)
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
		//获取演员获奖情况
		err, celebrity.Rewords = service.SelectRewards(celebrityId)
		if err != nil {
			fmt.Println("查询演员获奖情况错误", err)
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
		//获取演员最近五部作品
		err, celebrity.Works = service.SelectRecentWorks(celebrityId)
		if err != nil {
			fmt.Println("查询演员最近五部作品错误", err)
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
		celebrity.CelebrityId = celebrityId
		tool.ReturnSuccess(ctx, 200, celebrity)
		return
	}

	celebrityDetails := model.Celebrity{
		CelebrityId: celebrityId,
	}
	//姓名
	celebrityDetails.CelebrityName = docDetails.Find("#content > h1").Text()
	//头像
	image := docDetails.Find("#headline > div.pic > div > img")
	var ok bool
	celebrityDetails.Image, ok = image.Attr("src")
	if !ok {
		fmt.Println("演员图片解析失败", err)
		tool.ReturnFailure(ctx, 500, "图片解析失败")
		return
	}
	docDetails.Find("#headline > div.info > ul").
		Each(func(i int, selection *goquery.Selection) {
			gender := selection.Find("li:nth-child(1)").Text()
			constellation := selection.Find("li:nth-child(2)").Text()
			birthDate := selection.Find("li:nth-child(3)").Text()
			birthplace := selection.Find("li:nth-child(4)").Text()
			jobs := selection.Find("li:nth-child(5)").Text()
			nickname := selection.Find("li:nth-child(6)").Text()
			family := selection.Find("li:nth-child(7)").Text()
			//性别
			g, err := regexp.Compile(`[^\n\s*\r性别:(.*)\\n]`)
			if err != nil {
				fmt.Println("性别解析失败", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
				return
			}
			celebrityDetails.Gender = string(g.Find([]byte(gender)))
			//星座
			c, err := regexp.Compile(`[^\n\s*\r星座:\\n](.*)`)
			if err != nil {
				fmt.Println("性别解析失败", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
				return
			}
			celebrityDetails.Constellation = string(c.Find([]byte(constellation)))
			//出生地
			bPlace, err := regexp.Compile(`[^\n\s*\r出生地:\\n](.*)`)
			if err != nil {
				fmt.Println("出生地解析失败", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
			}
			celebrityDetails.Birthplace = string(bPlace.Find([]byte(birthplace)))
			//出生日期
			bDate, err := regexp.Compile(`[^\n\s*\r出生日期:\\n](.*)`)
			if err != nil {
				fmt.Println("出生日期解析错误", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
				return
			}
			celebrityDetails.BirthDate = string(bDate.Find([]byte(birthDate)))
			//职业
			j, err := regexp.Compile(`[^\n\s*\r职业:\\n](.*)`)
			if err != nil {
				fmt.Println("职业解析错误", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
				return
			}
			celebrityDetails.Jobs = string(j.Find([]byte(jobs)))
			//更多中文名
			n, err := regexp.Compile(`[^\n\s*\r更多中文名:\\n](.*)`)
			if err != nil {
				fmt.Println("更多中文名解析错误", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
				return
			}
			celebrityDetails.Nickname = string(n.Find([]byte(nickname)))
			//家庭成员
			f, err := regexp.Compile(`[^\n\s*\r家庭成员:\\n](.*)`)
			if err != nil {
				fmt.Println("家庭成员解析错误", err)
				tool.ReturnFailure(ctx, 500, "服务器错误")
				return
			}
			celebrityDetails.Family = string(f.Find([]byte(family)))
		})
	//演员简介
	celebrityDetails.Introduction = docDetails.Find("#intro > div.bd > span.all.hidden").Text()
	//存入演员基本信息
	err = service.InsertCelebrityDetails(celebrityDetails)
	if err != nil {
		fmt.Println("存入演员基本信息失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//照片
	for i := 1; i <= 5; i++ {
		j := strconv.Itoa(i)
		image = docDetails.Find("#photos > ul > li:nth-child(" + j + ") > a > img")
		picture, ok := image.Attr("src")
		if !ok {
			fmt.Println("图片解析失败")
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
		//存储照片
		err := service.InsertPhotos(picture, celebrityId)
		if err != nil {
			fmt.Println("照片存储失败", err)
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
	}
	/*
		#content > div > div.article > div:nth-child(5) > ul:nth-child(2) > li:nth-child(1)
		#content > div > div.article > div:nth-child(5) > ul:nth-child(2) > li:nth-child(2) > a
		#content > div > div.article > div:nth-child(5) > ul:nth-child(2) > li:nth-child(3)
		#content > div > div.article > div:nth-child(5) > ul:nth-child(2) > li:nth-child(4) > a
	*/
	//获奖情况
	for i := 2; i <= 4; i++ {
		j := strconv.Itoa(i)
		rewards := docDetails.Find("#content > div > div.article > div:nth-child(5) > ul:nth-child("+j+")> li:nth-child(1)").Text() + docDetails.Find("#content > div > div.article > div:nth-child(5) > ul:nth-child("+j+")> li:nth-child(2)").Text() + docDetails.Find("#content > div > div.article > div:nth-child(5) > ul:nth-child("+j+")> li:nth-child(3)").Text() + docDetails.Find("#content > div > div.article > div:nth-child(5) > ul:nth-child("+j+")> li:nth-child(4)").Text()
		//存储获奖情况
		err := service.InsertRewards(rewards, celebrityId)
		if err != nil {
			fmt.Println("储存获奖情况失败", err)
			tool.ReturnFailure(ctx, 500, "服务器错误")
			return
		}
	}
	//最近五部作品
	var works model.RecentWorks
	docDetails.Find("#recent_movies > div.bd > ul.list-s").
		Each(func(i int, selection *goquery.Selection) {
			for j := 1; j <= 5; j++ {
				k := strconv.Itoa(j)
				//作品id
				url := selection.Find("li:nth-child(" + k + ") > div.pic > a")
				workId, ok := url.Attr("href")
				if !ok {
					fmt.Println("作品id解析失败")
					tool.ReturnFailure(ctx, 500, "服务器错误")
					return
				}
				w, err := regexp.Compile(`\d+`)
				if err != nil {
					fmt.Println("抓取演职人员错误")
					tool.ReturnFailure(ctx, 500, "电影详情加载失败")
					return
				}
				works.WorkId = string(w.Find([]byte(workId)))
				//作品海报
				image := selection.Find("li:nth-child(" + k + ") > div.pic > a > img")
				imageTmp, ok := image.Attr("src")
				if !ok {
					fmt.Println("最近电影海报解析失败")
					tool.ReturnFailure(ctx, 500, "服务器错误")
					return
				}
				works.WorkImage = imageTmp
				//作品名字
				works.WorkName = selection.Find("li:nth-child(" + k + ") > div.info > a").Text()
				//作品评分
				works.WorkScores = selection.Find("li:nth-child(" + k + ") > div.info > em").Text()
				if works.WorkScores == "" {
					works.WorkScores = "暂无评分"
				}
				//储存最近五部作品
				err = service.InsertRecentWorks(works, celebrityId)
				if err != nil {
					fmt.Println("储存最近五部电影失败", err)
					tool.ReturnFailure(ctx, 500, "服务器错误")
					return
				}
			}
		})
	//获取演员信息
	err, celebrity := service.SelectCelebrityDetails(celebrityId)
	if err != nil {
		fmt.Println("查询演员信息错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//获取演员照片
	err, celebrity.Photos = service.SelectPhotos(celebrityId)
	if err != nil {
		fmt.Println("查询演员照片失败", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//获取演员获奖情况
	err, celebrity.Rewords = service.SelectRewards(celebrityId)
	if err != nil {
		fmt.Println("查询演员获奖情况错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	//获取演员最近五部作品
	err, celebrity.Works = service.SelectRecentWorks(celebrityId)
	if err != nil {
		fmt.Println("查询演员最近五部作品错误", err)
		tool.ReturnFailure(ctx, 500, "服务器错误")
		return
	}
	celebrity.CelebrityId = celebrityId
	tool.ReturnSuccess(ctx, 500, celebrity)
}
