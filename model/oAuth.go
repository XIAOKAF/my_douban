package model

import "gorm.io/gorm"

type OAuth struct {
	gorm.Model
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type Token struct {
	gorm.Model
	AccessToken string
}
