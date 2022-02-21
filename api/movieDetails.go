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
	"strings"
)

func selectMovieDetails(ctx *gin.Context) {
	tag := ctx.PostForm("tag")
	movieID := ctx.PostForm("movieID")

	//查找数据库内是否已存入该部电影，若未存入先insert
	err, flag := service.SelectMovieId(movieID)
	if err != nil {
		fmt.Println("查询失败", err)
		tool.ReturnFailure(ctx, 500, "加载电影详情失败")
		return
	}
	//该部电影已存入flag为false
	if flag == false {
		//获取除评分占比,同类比较，演职人员以外的信息
		err, movie := service.SelectMovieDetailsByMovieId(movieID)
		if err != nil {
			fmt.Println("查询失败", err)
			tool.ReturnFailure(ctx, 500, "电影详情加载失败")
			return
		}
		//获取评分占比
		err, starArr := service.SelectStarDetailsByMovieId(movieID)
		if err != nil {
			fmt.Println("获取评分占比失败", err)
			tool.ReturnFailure(ctx, 500, "电影详情失败")
			return
		}
		movie.StarPercentage = starArr
		//获取演职人员的信息
		err, actorArr := service.SelectAllActorsByMovieId(movieID)
		if err != nil {
			fmt.Println("获取演职人员基本信息失败", err)
			tool.ReturnFailure(ctx, 500, "电影详情加载失败")
			return
		}
		movie.AllRoles = actorArr
		tool.ReturnSuccess(ctx, 200, movie)
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/subject/"+movieID+"/?tag="+tag+"&from=gaia", nil)
	if err != nil {
		fmt.Println("请求失败", err)
		tool.ReturnFailure(ctx, 500, "加载失败")
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

	movieDetails := model.MovieDetails{
		MovieId: movieID,
	}

	movieDetails.MovieName = docDetails.Find("#content > h1 > span:nth-child(1)").Text()
	movieDetails.ReleaseYear = docDetails.Find("#content > h1 > span.year").Text()
	image := docDetails.Find("#mainpic > a > img")
	var ok bool
	movieDetails.Image, ok = image.Attr("src")
	if !ok {
		fmt.Println("图片解析失败")
		tool.ReturnFailure(ctx, 500, "电影详情加载失败")
		return
	}

	info := docDetails.Find("#info")
	for _, v := range strings.Split(info.Text(), "\n") {
		v = strings.TrimSpace(v)
		if strings.Contains(v, "导演") {
			movieDetails.Director = v
		}
		if strings.Contains(v, "编剧") {
			movieDetails.Author = v
		}
		if strings.Contains(v, "主演") {
			movieDetails.Actors = v
		}
		if strings.Contains(v, "类型") {
			movieDetails.Type = v
		}
		if strings.Contains(v, "制片国家/地区") {
			movieDetails.ProduceCountry = v
		}
		if strings.Contains(v, "语言") {
			movieDetails.Language = v
		}
		if strings.Contains(v, "上映日期") {
			movieDetails.ReleaseDate = v
		}
		if strings.Contains(v, "片长") {
			movieDetails.Duration = v
		}
		if strings.Contains(v, "又名") {
			movieDetails.Nickname = v
		}
	}

	var starPercentage model.StarPercentage
	docDetails.Find("#interest_sectl > div.rating_wrap.clearbox").
		Each(func(i int, selection *goquery.Selection) {
			movieDetails.RatingValue = selection.Find("div.rating_self.clearfix > strong").Text()
			ratingCount := selection.Find("div.rating_self.clearfix > div > div.rating_sum > a").Text()
			r, err := regexp.Compile(`\d+(.*)`)
			if err != nil {
				fmt.Println("评分人数解析失败", err)
				tool.ReturnFailure(ctx, 500, "电影详情加载失败")
				return
			}
			movieDetails.RatingCount = string(r.Find([]byte(ratingCount)))
			j := 1
			m := 5
			for j <= 5 && m >= 1 {
				k := strconv.Itoa(j)
				n := strconv.Itoa(m)
				star := selection.Find("div.ratings-on-weight > div:nth-child(" + k + ") > span.stars" + n + ".starstop").Text()
				percentage := selection.Find("div.ratings-on-weight > div:nth-child(" + k + ") > span.rating_per").Text()
				s, err := regexp.Compile(`\d+(.*)`)
				if err != nil {
					fmt.Println("星星解析失败", err)
					tool.ReturnFailure(ctx, 500, "电影详情加载失败")
					return
				}
				p, err := regexp.Compile(`\d+(.*)`)
				if err != nil {
					fmt.Println("星星占比解析失败", err)
					tool.ReturnFailure(ctx, 500, "电影详情加载失败")
					return
				}
				starPercentage.Star = string(s.Find([]byte(star)))
				starPercentage.Percentage = string(p.Find([]byte(percentage)))
				//插入评分占比
				err = service.InsertMovieScoresDetails(movieID, starPercentage)
				if err != nil {
					fmt.Println("评分详情插入失败", err)
					tool.ReturnFailure(ctx, 500, "电影详情加载失败")
					return
				}
				j++
				m--
			}
		})

	compare := "好于" + docDetails.Find("#interest_sectl > div.rating_betterthan > a:nth-child(1)").Text() + "  " + "好于" + docDetails.Find("#interest_sectl > div.rating_betterthan > a:nth-child(3)").Text()
	movieDetails.Compare = compare

	//电影简介
	description := docDetails.Find("#link-report > span").Text()
	d, err := regexp.Compile(`\S+(.*)`)
	if err != nil {
		fmt.Println("剧情简介解析失败", err)
		tool.ReturnFailure(ctx, 500, "电影详情加载失败")
		return
	}
	movieDetails.Description = string(d.Find([]byte(description)))

	var actorsBasicInfo model.ActorsBasicInfo
	for p := 1; p <= 6; p++ {
		q := strconv.Itoa(p)
		//演职人员id
		Url := docDetails.Find("#celebrities > ul > li:nth-child(" + q + ") > a")
		url, ok := Url.Attr("href")
		if !ok {
			fmt.Println("抓取演职人员错误")
			tool.ReturnFailure(ctx, 500, "电影详情加载失败")
			return
		}
		u, err := regexp.Compile(`\d+`)
		if err != nil {
			fmt.Println("抓取演职人员错误")
			tool.ReturnFailure(ctx, 500, "电影详情加载失败")
			return
		}
		actorsBasicInfo.CelebrityId = string(u.Find([]byte(url)))
		//演职人员image
		i := docDetails.Find("#celebrities > ul > li:nth-child(" + q + ") > a > div")
		im, ok := i.Attr("style")
		if !ok {
			fmt.Println(ok)
			return
		}
		s, err := regexp.Compile(`[a-zA-z]+://[^\s]*`)
		actorsBasicInfo.CelebrityImage = string(s.Find([]byte(im)))
		//演职人员姓名
		actorsBasicInfo.CelebrityName = docDetails.Find("#celebrities > ul > li:nth-child(" + q + ") > div > span.name > a").Text()
		//演职人员角色
		actorsBasicInfo.Role = docDetails.Find("#celebrities > ul > li:nth-child(" + q + ") > div > span.role").Text()
		//插入演职人员信息
		err = service.InsertActorBasis(movieID, actorsBasicInfo)
		if err != nil {
			fmt.Println("插入演职人员信息失败", err)
			tool.ReturnFailure(ctx, 500, "电影详情加载失败")
			return
		}
	}
	//插入除评分占比,演职人员以外的信息
	err = service.InsertMovieDetails(movieDetails)
	if err != nil {
		fmt.Println("影片信息插入失败", err)
		tool.ReturnFailure(ctx, 500, "影片详情加载失败")
		return
	}
	//获取除评分占比,同类比较，演职人员以外的信息
	err, movie := service.SelectMovieDetailsByMovieId(movieID)
	if err != nil {
		fmt.Println("查询失败", err)
		tool.ReturnFailure(ctx, 500, "电影详情加载失败")
		return
	}
	//获取评分占比
	err, starArr := service.SelectStarDetailsByMovieId(movieID)
	if err != nil {
		fmt.Println("获取评分占比失败", err)
		tool.ReturnFailure(ctx, 500, "电影详情失败")
		return
	}
	movie.StarPercentage = starArr
	//获取演职人员的信息
	err, actorArr := service.SelectAllActorsByMovieId(movieID)
	if err != nil {
		fmt.Println("获取演职人员基本信息失败", err)
		tool.ReturnFailure(ctx, 500, "电影详情加载失败")
		return
	}
	movie.AllRoles = actorArr
	tool.ReturnSuccess(ctx, 200, movie)
}

func PostShortComment(ctx *gin.Context) {
	movieId := ctx.PostForm("movieId")
	content := ctx.PostForm("comment")
	token := ctx.GetHeader("token")
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
	//token过期
	flag = tool.CheckToken(ctx, err)
	if !flag {
		return
	}
	if movieId == "" {
		tool.ReturnFailure(ctx, 403, "电影id不能为空")
		return
	}
	if content == "" {
		tool.ReturnFailure(ctx, 403, "欢迎一切有用的评论")
		return
	}
	if len(content) > 50 {
		tool.ReturnFailure(ctx, 403, "评论内容过多")
		return
	}
	comment := model.Comment{
		MovieId: movieId,
		Comment: content,
		PostId:  claims.Mobile,
	}
	//储存
	err = service.InsertComment(comment)
	if err != nil {
		fmt.Println("保存评论失败", err)
		tool.ReturnFailure(ctx, 500, "保存评论失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, "发布评论成功")
}
