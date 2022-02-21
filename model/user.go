package model

import "time"

type User struct {
	UserId           int
	Username         string
	Mobile           string
	Password         string
	VerifyCode       string
	SendTime         time.Time
	SelfIntroduction string
}
