package service

import (
	"database/sql"
	"gin/dao"
	"gin/model"
)

func SelectCelebrityById(celebrity string) (error, bool) {
	err := dao.SelectCelebrityById(celebrity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		return err, true
	}
	return nil, true
}

func InsertCelebrityDetails(celebrity model.Celebrity) error {
	err := dao.InsertCelebrityDetails(celebrity)
	if err != nil {
		return err
	}
	return nil
}

func InsertPhotos(picture string, celebrityId string) error {
	err := dao.InsertPhotos(picture, celebrityId)
	if err != nil {
		return err
	}
	return nil
}

func InsertRewards(rewards string, celebrityId string) error {
	err := dao.InsertRewards(rewards, celebrityId)
	if err != nil {
		return err
	}
	return nil
}

func InsertRecentWorks(works model.RecentWorks, celebrityId string) error {
	err := dao.InsertRecentWorks(works, celebrityId)
	if err != nil {
		return err
	}
	return nil
}

func SelectCelebrityDetails(celebrityId string) (error, model.Celebrity) {
	err, celebrity := dao.SelectCelebrityDetails(celebrityId)
	if err != nil {
		return err, celebrity
	}
	return nil, celebrity
}

func SelectPhotos(celebrityId string) (error, [5]string) {
	err, photoArr := dao.SelectPhotos(celebrityId)
	if err != nil {
		return err, photoArr
	}
	return nil, photoArr
}

func SelectRewards(celebrityId string) (error, [3]string) {
	err, rewardsArr := dao.SelectRewards(celebrityId)
	if err != nil {
		return err, rewardsArr
	}
	return nil, rewardsArr
}

func SelectRecentWorks(celebrityId string) (error, [5]model.RecentWorks) {
	err, worksArr := dao.SelectRecentWorks(celebrityId)
	if err != nil {
		return err, worksArr
	}
	return nil, worksArr
}
