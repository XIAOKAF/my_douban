package model

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"time"
)

type TokenClaims struct {
	gorm.Model
	Mobile     string
	Variety    string //token与refreshToken两种类型
	ExpireTime time.Time
	jwt.StandardClaims
}
