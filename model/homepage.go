package model

type HotShowing struct {
	Rank        int
	MovieName   string
	RatingValue string
	Image       string
}

type RecentHot struct {
	Rank             int
	RecentHotMovieId string
	MovieName        string
	RatingValue      string
	Image            string
}

type RecentHotTeleplay struct {
	Rank                int
	RecentHotTeleplayId string
	TeleplayName        string
	Update              string
	RatingValue         string
	Image               string
}

type WeeklyPraise struct {
	Rank                  int
	WeeklyPraiseMovieName string
}

type HotRecommendation struct {
	Rank    int
	Title   string
	Content string
	Image   string
}
