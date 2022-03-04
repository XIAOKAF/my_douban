package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenClaims struct {
	Mobile     string
	Variety    string //token与refreshToken两种类型
	ExpireTime time.Time
	jwt.StandardClaims
}
