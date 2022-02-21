package service

import (
	"gin/dao"
	"gin/model"
)

func TruncateInfo(table string) error {
	err := dao.TruncateInfo(table)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHotShowing(hotShowing model.HotShowing) error {
	err := dao.InsertHotShowing(hotShowing)
	if err != nil {
		return err
	}
	return nil
}

func SelectHotShowing(rank int) (model.HotShowing, error) {
	movieName, ratingValue, image, r, err := dao.SelectHotShowing(rank)
	hotShowing := model.HotShowing{
		Rank:        r,
		MovieName:   movieName,
		RatingValue: ratingValue,
		Image:       image,
	}
	if err != nil {
		return hotShowing, err
	}
	return hotShowing, nil
}

func UpdateRecentHotMovie(recentHotMovie model.RecentHot) error {
	err := dao.UpdateRecentHotMovie(recentHotMovie)
	if err != nil {
		return err
	}
	return nil
}

func SelectRecentHotMovie() ([50]model.RecentHot, error) {
	arr, err := dao.SelectRecentHotMovie()
	if err != nil {
		return arr, err
	}

	return arr, nil
}

func UpdateRecentHotTeleplay(recentHotTeleplay model.RecentHotTeleplay) error {
	err := dao.UpdateRecentHotTeleplay(recentHotTeleplay)
	if err != nil {
		return err
	}
	return nil
}

func SelectRecentHotTeleplay() ([50]model.RecentHotTeleplay, error) {
	arr, err := dao.SelectRecentHotTeleplay()
	if err != nil {
		return arr, err
	}
	return arr, nil
}

func UpdateWeeklyPraise(weeklyPraiseMovieName string) error {
	err := dao.UpdateWeeklyPraise(weeklyPraiseMovieName)
	if err != nil {
		return err
	}
	return nil
}

func SelectWeeklyPraise() ([10]model.WeeklyPraise, error) {
	arr, err := dao.SelectWeeklyPraise()
	if err != nil {
		return arr, err
	}
	return arr, nil
}

func UpdateHotRecommendation(hotRecommendation model.HotRecommendation) error {
	err := dao.UpdateHotRecommendation(hotRecommendation)
	if err != nil {
		return err
	}
	return nil
}

func SelectHotRecommendation() ([8]model.HotRecommendation, error) {
	arr, err := dao.SelectHotRecommendation()
	if err != nil {
		return arr, err
	}
	return arr, nil
}

func SelectMoviesByKeyWords(keyWords string) (error, [10]string) {
	err, movieArr := dao.SelectMovieByKeyWords(keyWords)
	if err != nil {
		return err, movieArr
	}
	return nil, movieArr
}
