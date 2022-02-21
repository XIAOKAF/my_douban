package api

import (
	"encoding/json"
	"fmt"
	"gin/model"
	"gin/service"
	"gin/tool"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

//正在热映
func hotShowing(ctx *gin.Context) {
	//发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/", nil)
	if err != nil {
		fmt.Println("请求错误：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}

	//加入一些请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.76")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "ll=\"108309\"; bid=gtXWf_veA68; push_doumail_num=0; push_noty_num=0; __utmv=30149280.21338; __yadk_uid=ushdXmeu6hPz8VL7FLqtxfB7XyoIDI0K; __gads=ID=077fbb5baee62c93-227b63397dd000ae:T=1644204395:RT=1644204395:S=ALNI_MawllVEXeybHBZ7lVBDC1sUfvJgCg; ct=y; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; dbcl2=\"213387422:/wRfw1lR4Aw\"; Hm_lvt_16a14f3002af32bf3a75dfe352478639=1644762356; ck=O7li; ap_v=0,6.0; __utmc=30149280; __utmc=223695111; _pk_ref.100001.4cf6=[\"\",\"\",1644938033,\"https://www.douban.com/\"]; _pk_ses.100001.4cf6=*; __utma=30149280.393236484.1643783658.1644933929.1644938033.33; __utmz=30149280.1644938033.33.18.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utma=223695111.2053381400.1643783658.1644933933.1644938033.31; __utmb=223695111.0.10.1644938033; __utmz=223695111.1644938033.31.14.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmb=30149280.2.10.1644938033; _pk_id.100001.4cf6=63f53dbc9e9bd5ef.1643783657.31.1644938380.1644933933.")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	defer resp.Body.Close()

	//解析网页
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析错误：", err)
		tool.ReturnFailure(ctx, 500, "网页解析错误")
		return
	}

	err = service.TruncateInfo("hotShowing")
	if err != nil {
		fmt.Println("删除最近热映电影失败", err)
		tool.ReturnFailure(ctx, 500, "正在热映加载失败")
		return
	}

	//定义存放数据的数组
	var arr [19]model.HotShowing

	//根据规律获取数据
	for j := 1; j <= 19; j++ {
		k := strconv.Itoa(j)
		docDetails.Find("#screening > div.screening-bd > ul> li:nth-child(" + k + ") > ul").
			Each(func(i int, selection *goquery.Selection) {
				movieName := selection.Find("li.title > a").Text()
				score := selection.Find("li.rating > span.subject-rate").Text()
				image := selection.Find("li.poster > a > img")
				imgTmp, ok := image.Attr("src")
				if score == "" {
					score = "暂无评分"
				}
				if !ok {
					fmt.Println("电影海报解析错误")
					tool.ReturnFailure(ctx, 500, "海报加载失败")
					return
				}
				movie := model.HotShowing{
					Rank:        j,
					MovieName:   movieName,
					RatingValue: score,
					Image:       imgTmp,
				}
				//更新正在热映影片消息
				err := service.UpdateHotShowing(movie)
				if err != nil {
					fmt.Println("正在热映数据插入失败", err)
					tool.ReturnFailure(ctx, 500, "正在热映加载失败")
					return
				}
				//获取正在热映影片消息
				m, err := service.SelectHotShowing(j)
				if err != nil {
					fmt.Println("获取正在热映数据失败", err)
					tool.ReturnFailure(ctx, 500, "正在热映加载失败")
					return
				}
				//将影片数据以结构体的形式存入数组
				arr[j-1] = m
			})
	}
	tool.ReturnSuccess(ctx, 200, arr)
}

//最近热门电影
func recentHotMovie(ctx *gin.Context) {
	choice := ctx.PostForm("choice")
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/j/search_subjects?type=movie&tag="+choice+"&page_limit=50&page_start=0", nil)
	if err != nil {
		fmt.Println("请求错误：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.76")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "ll=\"108309\"; bid=gtXWf_veA68; push_doumail_num=0; push_noty_num=0; __utmv=30149280.21338; __yadk_uid=ushdXmeu6hPz8VL7FLqtxfB7XyoIDI0K; __gads=ID=077fbb5baee62c93-227b63397dd000ae:T=1644204395:RT=1644204395:S=ALNI_MawllVEXeybHBZ7lVBDC1sUfvJgCg; ct=y; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; dbcl2=\"213387422:/wRfw1lR4Aw\"; Hm_lvt_16a14f3002af32bf3a75dfe352478639=1644762356; ck=O7li; ap_v=0,6.0; __utmc=30149280; __utmc=223695111; _pk_ref.100001.4cf6=[\"\",\"\",1644938033,\"https://www.douban.com/\"]; _pk_ses.100001.4cf6=*; __utma=30149280.393236484.1643783658.1644933929.1644938033.33; __utmz=30149280.1644938033.33.18.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utma=223695111.2053381400.1643783658.1644933933.1644938033.31; __utmb=223695111.0.10.1644938033; __utmz=223695111.1644938033.31.14.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmb=30149280.2.10.1644938033; _pk_id.100001.4cf6=63f53dbc9e9bd5ef.1643783657.31.1644938380.1644933933.")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	defer resp.Body.Close()
	recentHotMovie := model.RecentHot{}
	//存放所有从json里面获取的数据
	var data map[string]interface{}
	//获取json
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("获取网络json失败", err)
		tool.ReturnFailure(ctx, 500, "最近热门电影加载失败")
		return
	}
	//将json里面的数据映射到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("json数据映射失败", err)
		tool.ReturnFailure(ctx, 500, "最近热门电影加载失败")
		return
	}
	//将data中的subjects部分转化为[]interface{}类型
	movieInfo := data["subjects"].([]interface{})
	//删除原来的数据
	err = service.TruncateInfo("recentHotMovie")
	if err != nil {
		fmt.Println("删除最近热门电影失败", err)
		tool.ReturnFailure(ctx, 500, "最近热门电影加载失败")
		return
	}
	//循环便利将数据存入结构体
	for _, item := range movieInfo {
		movie := item.(map[string]interface{})
		recentHotMovie.Image = movie["cover"].(string)
		recentHotMovie.RecentHotMovieId = movie["id"].(string)
		recentHotMovie.MovieName = movie["title"].(string)
		recentHotMovie.RatingValue = movie["rate"].(string)
		//将获取到的信息存入数据库
		err := service.UpdateRecentHotMovie(recentHotMovie)
		if err != nil {
			fmt.Println("更新最近热映电影失败", err)
			tool.ReturnFailure(ctx, 500, "最近热映电影加载失败")
			return
		}
	}
	//从数据库里面获取信息并返回给前端
	movie, err := service.SelectRecentHotMovie()
	if err != nil {
		fmt.Println("获取最近热映电影失败", err)
		tool.ReturnFailure(ctx, 500, "最近热映电影加载失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, movie)
}

//最近热门电视剧
func recentHotTeleplay(ctx *gin.Context) {
	tag := ctx.PostForm("tag")
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/j/search_subjects?type=tv&tag="+tag+"&page_limit=50&page_start=0", nil)
	if err != nil {
		fmt.Println("请求错误：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.76")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "ll=\"108309\"; bid=gtXWf_veA68; push_doumail_num=0; push_noty_num=0; __utmv=30149280.21338; __yadk_uid=ushdXmeu6hPz8VL7FLqtxfB7XyoIDI0K; __gads=ID=077fbb5baee62c93-227b63397dd000ae:T=1644204395:RT=1644204395:S=ALNI_MawllVEXeybHBZ7lVBDC1sUfvJgCg; ct=y; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; dbcl2=\"213387422:/wRfw1lR4Aw\"; Hm_lvt_16a14f3002af32bf3a75dfe352478639=1644762356; ck=O7li; ap_v=0,6.0; __utmc=30149280; __utmc=223695111; _pk_ref.100001.4cf6=[\"\",\"\",1644938033,\"https://www.douban.com/\"]; _pk_ses.100001.4cf6=*; __utma=30149280.393236484.1643783658.1644933929.1644938033.33; __utmz=30149280.1644938033.33.18.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utma=223695111.2053381400.1643783658.1644933933.1644938033.31; __utmb=223695111.0.10.1644938033; __utmz=223695111.1644938033.31.14.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmb=30149280.2.10.1644938033; _pk_id.100001.4cf6=63f53dbc9e9bd5ef.1643783657.31.1644938380.1644933933.")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	defer resp.Body.Close()
	//删除原有数据
	err = service.TruncateInfo("recentHotTeleplay")
	if err != nil {
		fmt.Println("删除最近热门电视剧失败")
		tool.ReturnFailure(ctx, 500, "最近热门电视剧加载失败")
		return
	}
	recentHotTeleplay := model.RecentHotTeleplay{}
	//存放所有从json里面获取的数据
	var data map[string]interface{}
	//获取json
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("获取网络json失败", err)
		tool.ReturnFailure(ctx, 500, "最近热门电视剧加载失败")
		return
	}
	//将json里面的数据映射到data中
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("json数据映射失败", err)
		tool.ReturnFailure(ctx, 500, "最近热门电视剧加载失败")
		return
	}
	//将data中的subjects部分转化为[]interface{}类型
	teleplayInfo := data["subjects"].([]interface{})
	//循环便利将数据存入结构体
	for _, item := range teleplayInfo {
		teleplay := item.(map[string]interface{})
		recentHotTeleplay.Image = teleplay["cover"].(string)
		recentHotTeleplay.RecentHotTeleplayId = teleplay["id"].(string)
		recentHotTeleplay.TeleplayName = teleplay["title"].(string)
		recentHotTeleplay.Update = teleplay["episodes_info"].(string)
		recentHotTeleplay.RatingValue = teleplay["rate"].(string)
		//将获取到的信息存入数据库
		err := service.UpdateRecentHotTeleplay(recentHotTeleplay)
		if err != nil {
			fmt.Println("更新最近热门电视剧失败", err)
			tool.ReturnFailure(ctx, 500, "最近热门电视剧加载失败")
			return
		}
	}
	//从数据库里面获取信息并返回给前端
	arr, err := service.SelectRecentHotTeleplay()
	if err != nil {
		fmt.Println("获取最近热门电视剧失败", err)
		tool.ReturnFailure(ctx, 500, "最近热门电视剧加载失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, arr)
}

//一周热榜
func weeklyPraise(ctx *gin.Context) {
	//1、发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/", nil)
	if err != nil {
		fmt.Println("请求错误：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	//加入一些请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.76")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	defer resp.Body.Close()

	//2、解析网页
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析错误：", err)
		tool.ReturnFailure(ctx, 500, "一周热榜加载失败")
		return
	}

	err = service.TruncateInfo("weeklyPraise")
	if err != nil {
		fmt.Println("删除一周热榜数据失败", err)
		tool.ReturnFailure(ctx, 500, "一周热榜加载失败")
		return
	}

	for j := 1; j <= 10; j++ {
		k := strconv.Itoa(j)
		weeklyPraiseMovieName := docDetails.Find("#billboard > div.billboard-bd > table > tbody > tr:nth-child(" + k + ") > td.title > a").Text()
		err := service.UpdateWeeklyPraise(weeklyPraiseMovieName)
		if err != nil {
			fmt.Println("一周热榜更新失败", err)
			tool.ReturnFailure(ctx, 500, "一周热榜数据加载失败")
			return
		}
	}
	arr, err := service.SelectWeeklyPraise()
	if err != nil {
		fmt.Println("一周热榜数据获取失败", err)
		tool.ReturnFailure(ctx, 500, "一周热榜数据加载失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, arr)
}

//热门推荐
func hotRecommendation(ctx *gin.Context) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/", nil)
	if err != nil {
		fmt.Println("请求错误：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}

	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 Edg/97.0.1072.76")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.douban.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "ll=\"108309\"; bid=gtXWf_veA68; push_noty_num=0; push_doumail_num=0; __utmv=30149280.21338; __yadk_uid=ushdXmeu6hPz8VL7FLqtxfB7XyoIDI0K; __gads=ID=077fbb5baee62c93-227b63397dd000ae:T=1644204395:RT=1644204395:S=ALNI_MawllVEXeybHBZ7lVBDC1sUfvJgCg; ct=y; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; _vwo_uuid_v2=DA7B878DAE010FF2EA9148B12AE746A54|3d548311f6308d509827284220612809; dbcl2=\"213387422:/wRfw1lR4Aw\"; Hm_lvt_16a14f3002af32bf3a75dfe352478639=1644762356; __utmz=30149280.1644938033.33.18.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ck=O7li; __utmc=30149280; __utmc=223695111; __utmz=223695111.1645000504.35.17.utmcsr=baidu|utmccn=(organic)|utmcmd=organic|utmctr=豆瓣; _pk_ref.100001.4cf6=[\"\",\"\",1645024639,\"https://www.baidu.com/s?wd=%E8%B1%86%E7%93%A3&rsv_spt=1&rsv_iqid=0x94b40b230009fceb&issp=1&f=8&rsv_bp=1&rsv_idx=2&ie=utf-8&tn=15007414_8_dg&rsv_enter=1&rsv_dl=tb&rsv_sug3=10&rsv_sug1=9&rsv_sug7=100&rsv_sug2=0&rsv_btype=i&inputT=3652&rsv_sug4=5876\"]; _pk_id.100001.4cf6=63f53dbc9e9bd5ef.1643783657.36.1645024639.1645000504.; ap_v=0,6.0; __utma=30149280.393236484.1643783658.1645000504.1645024643.38; __utma=223695111.2053381400.1643783658.1645000504.1645024643.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		tool.ReturnFailure(ctx, 500, "网络爬虫请求失败")
		return
	}
	defer resp.Body.Close()

	//2、解析网页
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析错误：", err)
		tool.ReturnFailure(ctx, 500, "热门推荐加载失败")
		return
	}

	err = service.TruncateInfo("hotRecommendation")
	if err != nil {
		fmt.Println("热门推荐数据删除失败", err)
		tool.ReturnFailure(ctx, 500, "热门推荐加载失败")
		return
	}

	for j := 1; j <= 8; j++ {
		k := strconv.Itoa(j)
		docDetails.Find("#hot-gallery > ul > li:nth-child(" + k + ") > div").
			Each(func(i int, selection *goquery.Selection) {
				image := selection.Find("a > img")
				imageTmp, ok := image.Attr("src")
				title := selection.Find("div > div.gallery-hd > a > h3").Text()
				content := selection.Find("div > div.gallery-bd > p").Text()
				//正则表达式去除content中的空格
				c, err := regexp.Compile(`[^\x00-\xff]+`)
				if err != nil {
					fmt.Println("热门推荐内容解析失败", err)
					tool.ReturnFailure(ctx, 500, "内热门推荐加载失败")
					return
				}
				content = string(c.Find([]byte(content)))
				if !ok {
					fmt.Println("图片解析失败")
					tool.ReturnFailure(ctx, 500, "热门推荐加载失败")
					return
				}
				hotRecommendation := model.HotRecommendation{
					Title:   title,
					Content: content,
					Image:   imageTmp,
				}
				err = service.UpdateHotRecommendation(hotRecommendation)
				if err != nil {
					fmt.Println("热门推荐数据更新失败", err)
					tool.ReturnFailure(ctx, 500, "热门推荐加载失败")
					return
				}
			})
	}
	arr, err := service.SelectHotRecommendation()
	if err != nil {
		fmt.Println("热门推荐数据获取失败", err)
		tool.ReturnFailure(ctx, 500, "热门推荐加载失败")
		return
	}
	tool.ReturnSuccess(ctx, 200, arr)
}

func SelectMovieByKeyWords(ctx *gin.Context) {
	keyWords := ctx.PostForm("keywords")
	err, movieArr := service.SelectMoviesByKeyWords(keyWords)
	if err != nil {
		fmt.Println("查询失败", err)
		tool.ReturnFailure(ctx, 500, "查询错误")
		return
	}
	tool.ReturnSuccess(ctx, 200, movieArr)
}
