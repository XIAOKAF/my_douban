package dao

import (
	"gin/model"
)

// SelectMovieId 查询该部电影是否已经在数据库中
func SelectMovieId(movieId string) error {
	var name string
	rows := DB.QueryRow("SELECT movieName FROM movieDetails WHERE movieId = ? ", movieId)
	if rows.Err() != nil {
		return rows.Err()
	}
	err := rows.Scan(&name)
	if err != nil {
		return err
	}
	return nil
}

// InsertMovieDetails 插入除评分详情以及演职人员以外的信息
func InsertMovieDetails(movieDetails model.MovieDetails) error {
	sql := "insert into movieDetails(movieId,movieName,releaseYear,movieImage,director,author,actors,ty,produceCountry,lan,releaseDate,duration,nickname,ratingValue,ratingCount,compare,description)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := DB.Exec(sql, movieDetails.MovieId, movieDetails.MovieName, movieDetails.ReleaseYear, movieDetails.Image, movieDetails.Director, movieDetails.Author, movieDetails.Actors, movieDetails.Type, movieDetails.ProduceCountry, movieDetails.Language, movieDetails.ReleaseDate, movieDetails.Duration, movieDetails.Nickname, movieDetails.RatingValue, movieDetails.RatingCount, movieDetails.Compare, movieDetails.Description)
	if err != nil {
		return err
	}
	return nil
}

// InsertScoreDetails 插入评分详情
func InsertScoreDetails(worksId string, percentage model.StarPercentage) error {
	sql := "insert into star(worksID,star,per)values(?,?,?)"
	_, err := DB.Exec(sql, worksId, percentage.Star, percentage.Percentage)
	if err != nil {
		return err
	}
	return nil
}

// InsertActorsBasis 插入演职人员基本信息
func InsertActorsBasis(movieId string, info model.ActorsBasicInfo) error {
	sql := "insert into actorBasis(work,actorId,actorName,actorImage,actorRole)values(?,?,?,?,?)"
	_, err := DB.Exec(sql, movieId, info.CelebrityId, info.CelebrityName, info.CelebrityImage, info.Role)
	if err != nil {
		return err
	}
	return nil
}

// SelectMovieDetailsByMovieId 查询除评分详情以及演职人员以外的信息
func SelectMovieDetailsByMovieId(movieId string) (error, model.MovieDetails) {
	var movieDetails model.MovieDetails

	sql := "select movieId,movieName,releaseYear,movieImage,director,author,actors,ty,produceCountry,lan,releaseDate,duration,nickname,ratingValue,ratingCount,compare,description from movieDetails where movieId = ?"

	rows, err := DB.Query(sql, movieId)
	if err != nil {
		return err, movieDetails
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&movieDetails.MovieId, &movieDetails.MovieName, &movieDetails.ReleaseYear, &movieDetails.Image, &movieDetails.Director, &movieDetails.Author, &movieDetails.Actors, &movieDetails.Type, &movieDetails.ProduceCountry, &movieDetails.Language, &movieDetails.ReleaseDate, &movieDetails.Duration, &movieDetails.Nickname, &movieDetails.RatingValue, &movieDetails.RatingCount, &movieDetails.Compare, &movieDetails.Description)
		if err != nil {
			return err, movieDetails
		}
	}
	return nil, movieDetails
}

// SelectStarDetailsByMovieId 查询评分详情
func SelectStarDetailsByMovieId(movieId string) (error, [5]model.StarPercentage) {
	i := 0
	var arr [5]model.StarPercentage
	var per model.StarPercentage
	sql := "select star,per from star where worksId = ?"
	rows, err := DB.Query(sql, movieId)
	if err != nil {
		return err, arr
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&per.Star, &per.Percentage)
		if err != nil {
			return err, arr
		}
		arr[i] = per
		i++
	}
	return nil, arr
}

// SelectAllActorsByMovieId 查询演职人员基本信息
func SelectAllActorsByMovieId(movieId string) (error, [6]model.ActorsBasicInfo) {
	i := 0
	var arr [6]model.ActorsBasicInfo
	var basis model.ActorsBasicInfo
	sql := "select actorId,actorName,actorImage,actorRole from actorBasis where work = ?"
	rows, err := DB.Query(sql, movieId)
	if err != nil {
		return err, arr
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&basis.CelebrityId, &basis.CelebrityName, &basis.CelebrityImage, &basis.Role)
		if err != nil {
			return err, arr
		}
		arr[i] = basis
		i++
	}
	return nil, arr
}
