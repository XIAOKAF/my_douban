package model

import "gorm.io/gorm"

type HotShowing struct {
	gorm.Model
	Rank        int
	MovieName   string
	RatingValue string
	Image       string
}

type RecentHot struct {
	gorm.Model
	Rank             int
	RecentHotMovieId string
	MovieName        string
	RatingValue      string
	Image            string
}

type RecentHotTeleplay struct {
	gorm.Model
	Rank                int
	RecentHotTeleplayId string
	TeleplayName        string
	Update              string
	RatingValue         string
	Image               string
}

type WeeklyPraise struct {
	gorm.Model
	Rank                  int
	WeeklyPraiseMovieName string
}

type HotRecommendation struct {
	gorm.Model
	Rank    int
	Title   string
	Content string
	Image   string
}
