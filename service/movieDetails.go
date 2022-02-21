package service

import (
	"database/sql"
	"gin/dao"
	"gin/model"
)

func SelectMovieId(movieId string) (error, bool) {
	err := dao.SelectMovieId(movieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, true
		}
		return err, false
	}
	return nil, false
}

func InsertMovieDetails(movieDetails model.MovieDetails) error {
	err := dao.InsertMovieDetails(movieDetails)
	if err != nil {
		return err
	}
	return nil
}

func InsertMovieScoresDetails(worksId string, percentage model.StarPercentage) error {
	err := dao.InsertScoreDetails(worksId, percentage)
	if err != nil {
		return err
	}
	return nil
}

func InsertActorBasis(movieId string, info model.ActorsBasicInfo) error {
	err := dao.InsertActorsBasis(movieId, info)
	if err != nil {
		return err
	}
	return nil
}

func SelectMovieDetailsByMovieId(movieId string) (error, model.MovieDetails) {
	err, movieDetails := dao.SelectMovieDetailsByMovieId(movieId)
	if err != nil {
		return err, movieDetails
	}
	return nil, movieDetails
}

func SelectStarDetailsByMovieId(movieId string) (error, [5]model.StarPercentage) {
	err, arr := dao.SelectStarDetailsByMovieId(movieId)
	if err != nil {
		return err, arr
	}
	return nil, arr
}

func SelectAllActorsByMovieId(movieId string) (error, [6]model.ActorsBasicInfo) {
	err, arr := dao.SelectAllActorsByMovieId(movieId)
	if err != nil {
		return err, arr
	}
	return nil, arr
}
