package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserId           int
	Username         string
	Mobile           string
	Password         string
	VerifyCode       string
	SendTime         time.Time
	SelfIntroduction string
}
