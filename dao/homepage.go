package dao

import (
	"gin/model"
)

func TruncateInfo(table string) error {
	sql := "truncate table " + table
	_, err := DB.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func InsertHotShowing(hotShowing model.HotShowing) error {
	sql := "insert into hotShowing(movieName, ratingValue, image) values(?, ?, ?)"
	_, err := DB.Exec(sql, hotShowing.MovieName, hotShowing.RatingValue, hotShowing.Image)
	if err != nil {
		return err
	}
	return nil
}

func SelectHotShowing(rank int) (string, string, string, int, error) {
	var movieName string
	var ratingValue string
	var image string
	rows, err := DB.Query("select movieName, ratingValue, image from hotShowing where hotShowingId = ?", rank)
	if err != nil {
		return "", "", "", 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&movieName, &ratingValue, &image)
		if err != nil {
			return movieName, ratingValue, image, 0, err
		}
	}
	return movieName, ratingValue, image, rank, nil
}

func UpdateRecentHotMovie(recentHot model.RecentHot) error {
	sql := "insert into recentHotMovie(recentHotMovieId, movieName, ratingValue, image) values(?,?, ?, ?)"
	_, err := DB.Exec(sql, recentHot.RecentHotMovieId, recentHot.MovieName, recentHot.RatingValue, recentHot.Image)
	if err != nil {
		return err
	}
	return nil
}

func SelectRecentHotMovie() ([50]model.RecentHot, error) {
	var arr [50]model.RecentHot
	var recentHotMovie model.RecentHot
	rows, err := DB.Query("select r,recentHotMovieId,movieName,ratingValue,image from recentHotMovie")
	if err != nil {
		return arr, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&recentHotMovie.Rank, &recentHotMovie.RecentHotMovieId, &recentHotMovie.MovieName, &recentHotMovie.RatingValue, &recentHotMovie.Image)
		arr[recentHotMovie.Rank-1] = recentHotMovie
		if err != nil {
			return arr, err
		}
	}
	return arr, nil
}

func UpdateRecentHotTeleplay(recentHotTeleplay model.RecentHotTeleplay) error {
	sql := "insert into recentHotTeleplay(recentHotTeleplayId, teleplayName, updat, ratingValue, image)value(?,?,?,?,?)"
	_, err := DB.Exec(sql, recentHotTeleplay.RecentHotTeleplayId, recentHotTeleplay.TeleplayName, recentHotTeleplay.Update, recentHotTeleplay.RatingValue, recentHotTeleplay.Image)
	if err != nil {
		return err
	}
	return nil
}

func SelectRecentHotTeleplay() ([50]model.RecentHotTeleplay, error) {
	var arr [50]model.RecentHotTeleplay
	var recentHotTeleplay model.RecentHotTeleplay
	rows, err := DB.Query("select sort, recentHotTeleplayId, teleplayName, updat, ratingValue, image from recentHotTeleplay")
	if err != nil {
		return arr, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&recentHotTeleplay.Rank, &recentHotTeleplay.RecentHotTeleplayId, &recentHotTeleplay.TeleplayName, &recentHotTeleplay.Update, &recentHotTeleplay.RatingValue, &recentHotTeleplay.Image)
		arr[recentHotTeleplay.Rank-1] = recentHotTeleplay
		if err != nil {
			return arr, err
		}
	}
	return arr, nil
}

func UpdateWeeklyPraise(weeklyPraiseName string) error {
	sql := "insert into weeklyPraise (weeklyPraiseMovieName)value(?)"
	_, err := DB.Query(sql, weeklyPraiseName)
	if err != nil {
		return err
	}
	return nil
}

func SelectWeeklyPraise() ([10]model.WeeklyPraise, error) {
	var arr [10]model.WeeklyPraise
	var weeklyPraiseMovie model.WeeklyPraise
	rows, err := DB.Query("select s, weeklyPraiseMovieName from weeklyPraise")
	if err != nil {
		return arr, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&weeklyPraiseMovie.Rank, &weeklyPraiseMovie.WeeklyPraiseMovieName)
		if err != nil {
			return arr, err
		}
		arr[weeklyPraiseMovie.Rank-1] = weeklyPraiseMovie
	}
	return arr, nil
}

func UpdateHotRecommendation(hotRecommendation model.HotRecommendation) error {
	sql := "insert into hotRecommendation(title, content, image)values(?,?,?)"
	_, err := DB.Exec(sql, hotRecommendation.Title, hotRecommendation.Content, hotRecommendation.Image)
	if err != nil {
		return err
	}
	return nil
}

func SelectHotRecommendation() ([8]model.HotRecommendation, error) {
	var arr [8]model.HotRecommendation
	var hotRecommendation model.HotRecommendation
	rows, err := DB.Query("select ranking, title, content, image from hotRecommendation")
	if err != nil {
		return arr, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&hotRecommendation.Rank, &hotRecommendation.Title, &hotRecommendation.Content, &hotRecommendation.Image)
		if err != nil {
			return arr, err
		}
		arr[hotRecommendation.Rank-1] = hotRecommendation
	}
	return arr, nil
}

func SelectMovieByKeyWords(keyWords string) (error, [10]string) {
	var movieArr [10]string
	var i int
	sql := "select movieName from movieDetails where movieName like '%" + keyWords + "%'"
	rows, err := DB.Query(sql)
	if err != nil {
		return err, movieArr
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&movieArr[i])
		if err != nil {
			return err, movieArr
		}
		i++
	}
	return nil, movieArr
}
