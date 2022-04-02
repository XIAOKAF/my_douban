package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	MovieId string
	Comment string
	PostId  string
}
