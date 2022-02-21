package service

import (
	"errors"
	"gin/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("www.zxc.com")

func CreatToken(mobile string, variety string, duration time.Duration) (string, error) {
	expireTime := time.Now().Add(duration * time.Minute)
	claims := model.TokenClaims{
		Mobile:  mobile,
		Variety: variety,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(), //颁发时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*model.TokenClaims, error) {
	var tokenClaims model.TokenClaims
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, errors.New("token is expired")
	}
	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		return nil, errors.New("token解析失败")
	}
	err = token.Claims.Valid()
	if err != nil {
		return nil, err
	}
	return claims, nil
}
