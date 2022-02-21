package model

import (
	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	Mobile  string
	Variety string //token与refreshToken两种类型
	jwt.StandardClaims
}
