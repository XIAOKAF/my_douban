package model

import "gorm.io/gorm"

type MovieDetails struct {
	gorm.Model
	MovieId        string
	MovieName      string
	ReleaseYear    string
	Image          string
	Director       string
	Author         string
	Actors         string
	Type           string
	ProduceCountry string
	Language       string
	ReleaseDate    string
	Duration       string
	Nickname       string
	RatingValue    string
	RatingCount    string
	StarPercentage [5]StarPercentage
	Compare        string
	Description    string
	AllRoles       [6]ActorsBasicInfo
}

type StarPercentage struct {
	Star       string
	Percentage string
}

type ActorsBasicInfo struct {
	CelebrityId    string
	CelebrityImage string
	CelebrityName  string
	Role           string
}
