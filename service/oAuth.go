package service

import (
	"encoding/json"
	"fmt"
	"gin/model"
	"net/http"
)

func GetUrl(oAuth model.OAuth, code string) string {
	return fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", oAuth.ClientId, oAuth.ClientSecret, code)
}

func GetToken(url string, token model.Token) (*model.Token, error) {
	var req *http.Request
	var httpClient http.Client
	//形成请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &token, err
	}
	req.Header.Set("accept", "application/json")
	//发送请求并获得响应
	res, err := httpClient.Do(req)
	if err != nil {
		return &token, err
	}
	//将响应体解析为token并返回
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return &token, err
	}
	return &token, nil
}

func GetUserInfo(token *model.Token, user model.User) (model.User, error) {
	//形成请求
	userInfoUrl := "https://api.github.com/user"
	req, err := http.NewRequest(http.MethodGet, userInfoUrl, nil)
	if err != nil {
		return user, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))
	//发送请求并获取响应
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return user, err
	}
	//将数据写入user中
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CheckRegisterByUsername(user model.User) error {

	return nil
}

func RegisterByOAuth(user model.User) error {

	return nil
}
